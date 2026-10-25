package apps

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/texm/dokku-go"
	"gitlab.com/texm/shokku/internal/env"
	"gitlab.com/texm/shokku/internal/server/dto"
	"net/http"
	"regexp"
)

func isMountInSlice(mount dokku.StorageBindMount, mounts []dokku.StorageBindMount) bool {
	for _, m := range mounts {
		if m == mount {
			return true
		}
	}
	return false
}

func GetAppStorage(e *env.Env, c echo.Context) error {
	var req dto.GetAppStorageRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	report, err := e.Dokku.GetAppStorageReport(req.Name)
	if err != nil {
		return fmt.Errorf("getting app storage report: %w", err)
	}

	var allMounts []dokku.StorageBindMount
	allMounts = append(allMounts, report.RunMounts...)
	allMounts = append(allMounts, report.DeployMounts...)
	allMounts = append(allMounts, report.BuildMounts...)

	seenMap := map[string]bool{}
	mounts := []dto.StorageMount{}
	for _, dokkuMount := range allMounts {
		if _, seen := seenMap[dokkuMount.String()]; seen {
			continue
		}
		seenMap[dokkuMount.String()] = true
		mounts = append(mounts, dto.StorageMount{
			HostDir:       dokkuMount.HostDir,
			ContainerDir:  dokkuMount.ContainerDir,
			IsBuildMount:  isMountInSlice(dokkuMount, report.BuildMounts),
			IsRunMount:    isMountInSlice(dokkuMount, report.RunMounts),
			IsDeployMount: isMountInSlice(dokkuMount, report.DeployMounts),
		})
	}

	return c.JSON(http.StatusOK, dto.GetAppStorageResponse{
		Mounts: mounts,
	})
}

// dokkuStorageRegex matches a Dokku storage directory (/var/lib/dokku/data/storage/<name>).
// Anything else is a docker volume name or an absolute host path and is passed
// through to dokku, which validates it.
var dokkuStorageRegex = regexp.MustCompile(`^/var/lib/dokku/data/storage/(\w{2,})$`)

func MountAppStorage(e *env.Env, c echo.Context) error {
	var req dto.AlterAppStorageRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	// Ensure the storage directory exists when it is a Dokku storage path
	if matches := dokkuStorageRegex.FindStringSubmatch(req.HostDir); matches != nil {
		err := e.Dokku.EnsureStorageDirectory(matches[1], dokku.StorageChownOptionHerokuish)
		if err != nil {
			return fmt.Errorf("ensuring app storage dir: %w", err)
		}
	}

	mount := dokku.StorageBindMount{
		HostDir:      req.HostDir,
		ContainerDir: req.ContainerDir,
	}
	if err := e.Dokku.MountAppStorage(req.Name, mount); err != nil {
		return fmt.Errorf("mounting app storage dir: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func UnmountAppStorage(e *env.Env, c echo.Context) error {
	var req dto.AlterAppStorageRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	mount := dokku.StorageBindMount{
		HostDir:      req.HostDir,
		ContainerDir: req.ContainerDir,
	}
	if err := e.Dokku.UnmountAppStorage(req.Name, mount); err != nil {
		return fmt.Errorf("unmounting app storage dir: %w", err)
	}

	if req.RestartApp {
		go func() {
			if _, err := e.Dokku.RestartApp(req.Name, nil); err != nil {
				log.Error().Err(err).Msg("error while restarting app")
			}
		}()
	}

	return c.NoContent(http.StatusOK)
}
