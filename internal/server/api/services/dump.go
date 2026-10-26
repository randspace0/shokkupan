package services

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/dto"
)

// DownloadServiceDump streams `<type>:export <service>` to the browser, so a
// database can be backed up without configuring S3. The dump is never held in
// memory or written to disk on the Shokkupan side.
func DownloadServiceDump(e *env.Env, c echo.Context) error {
	var req dto.DownloadServiceDumpRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	dbSvc, err := lookupDBServiceByName(e, req.Name)
	if err != nil {
		return echo.ErrNotFound
	}

	// type comes from the database (redis/postgres/...), name was validated
	stream, err := e.Dokku.ExecStreaming(fmt.Sprintf("%s:export %s", dbSvc.Type, dbSvc.Name))
	if err != nil {
		return fmt.Errorf("starting export: %w", err)
	}

	filename := fmt.Sprintf("%s-%s.dump", dbSvc.Name, time.Now().UTC().Format("20060102-150405"))
	return streamDump(c, stream.Stdout, stream.Stderr, func() error { return stream.Error }, filename)
}

// streamDump copies out to the response. It waits for the first byte before
// sending headers so an immediate failure (unknown service, plugin missing)
// becomes a clear error instead of an empty "successful" download.
func streamDump(c echo.Context, out, errOut io.Reader, exitErr func() error, filename string) error {
	br := bufio.NewReader(out)
	if _, err := br.Peek(1); err != nil {
		msg, _ := io.ReadAll(io.LimitReader(errOut, 4096))
		return echo.NewHTTPError(http.StatusBadGateway, "export produced no data: "+string(msg))
	}

	h := c.Response().Header()
	h.Set(echo.HeaderContentType, "application/octet-stream")
	h.Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", filename))
	h.Set("Cache-Control", "no-store")
	c.Response().WriteHeader(http.StatusOK)

	if _, err := io.Copy(c.Response(), br); err != nil {
		return err
	}
	// headers are gone; a failure now can only be logged by the caller
	return exitErr()
}
