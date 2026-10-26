package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/access"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
)

const (
	roleContextKey = "user-role"
	appsContextKey = "user-apps"
	maxBodyForApp  = 1 << 20
)

// requestApp finds the app a request acts on: ?name= / ?app=, then the same
// keys in a JSON body (which is put back for the handler), then a form field.
func requestApp(c echo.Context) string {
	r := c.Request()
	for _, k := range []string{"name", "app"} {
		if v := c.QueryParam(k); v != "" {
			return v
		}
	}
	ct := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyForApp))
		if err != nil {
			return ""
		}
		r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), r.Body))
		var fields map[string]any
		if json.Unmarshal(body, &fields) != nil {
			return ""
		}
		for _, k := range []string{"name", "app"} {
			if v, ok := fields[k].(string); ok && v != "" {
				return v
			}
		}
	case strings.HasPrefix(ct, "multipart/form-data"), strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		return c.FormValue("name")
	}
	return ""
}

// loadAccess puts the caller's role and assigned apps on the context. A user
// that no longer exists gets no role, which Allows() treats as no access.
func loadAccess(e *env.Env, c echo.Context) (role string, apps map[string]bool) {
	user, err := e.Auth.GetUserFromContext(c)
	if err != nil {
		return "", nil
	}
	var u models.User
	if err := e.DB.Where("name = ?", user.Name).Take(&u).Error; err != nil {
		return "", nil
	}
	apps = map[string]bool{}
	if u.Role != access.RoleAdmin {
		var rows []models.AppAccess
		e.DB.Where("user_name = ?", u.Name).Find(&rows)
		for _, r := range rows {
			apps[r.AppName] = true
		}
	}
	c.Set(roleContextKey, u.Role)
	c.Set(appsContextKey, apps)
	return u.Role, apps
}

// Authorize enforces roles and per-app access on the protected API.
func Authorize(e *env.Env) echo.MiddlewareFunc {
	logger := middlewareLogger("authorize")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipDuringSetup(e, c) {
				return next(c)
			}
			rule := access.Classify(c.Request().Method, c.Request().URL.Path)
			role, apps := loadAccess(e, c)

			if !access.Allows(role, rule.Perm) {
				logger.Debug().Str("role", role).Str("path", c.Request().URL.Path).Msg("forbidden")
				return echo.NewHTTPError(http.StatusForbidden, "your role does not allow this")
			}
			if rule.Scoped && role != access.RoleAdmin {
				app := requestApp(c)
				if app == "" || !apps[app] {
					return echo.NewHTTPError(http.StatusForbidden, "you do not have access to this app")
				}
			}
			return next(c)
		}
	}
}

// VisibleApps returns the apps the caller may see; all is true for admins.
func VisibleApps(c echo.Context) (apps map[string]bool, all bool) {
	if role, _ := c.Get(roleContextKey).(string); role == access.RoleAdmin {
		return nil, true
	}
	apps, _ = c.Get(appsContextKey).(map[string]bool)
	return apps, false
}

// FilterApps keeps only the app names the caller may see.
func FilterApps(c echo.Context, names []string) []string {
	apps, all := VisibleApps(c)
	if all {
		return names
	}
	out := make([]string, 0, len(names))
	for _, n := range names {
		if apps[n] {
			out = append(out, n)
		}
	}
	return out
}
