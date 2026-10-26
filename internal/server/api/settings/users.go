package settings

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/access"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/randspace0/shokkupan/internal/server/dto"
	"gorm.io/gorm"
)

// SetUserRole changes a user's role. The last admin cannot be demoted, or
// nobody could manage users any more.
func SetUserRole(e *env.Env, c echo.Context) error {
	var req dto.SetUserRoleRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	err := e.DB.Transaction(func(tx *gorm.DB) error {
		var u models.User
		if err := tx.Where("name = ?", req.Name).Take(&u).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "no such user")
		}
		if u.Role == access.RoleAdmin && req.Role != access.RoleAdmin {
			var admins int64
			tx.Model(&models.User{}).Where("role = ?", access.RoleAdmin).Count(&admins)
			if admins <= 1 {
				return echo.NewHTTPError(http.StatusConflict, "cannot demote the last admin")
			}
		}
		return tx.Model(&u).Update("role", req.Role).Error
	})
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// SetUserApps replaces the set of apps a non-admin user can access.
func SetUserApps(e *env.Env, c echo.Context) error {
	var req dto.SetUserAppsRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	err := e.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&models.User{}).Where("name = ?", req.Name).Count(&count)
		if count == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "no such user")
		}
		if err := tx.Unscoped().Where("user_name = ?", req.Name).Delete(&models.AppAccess{}).Error; err != nil {
			return fmt.Errorf("clearing app access: %w", err)
		}
		seen := map[string]bool{}
		for _, app := range req.Apps {
			if seen[app] {
				continue
			}
			seen[app] = true
			if err := tx.Create(&models.AppAccess{UserName: req.Name, AppName: app}).Error; err != nil {
				return fmt.Errorf("granting %s: %w", app, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
