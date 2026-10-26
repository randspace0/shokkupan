package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/randspace0/shokkupan/internal/server/auth"
)

// stubAuth returns a fixed user from the context.
type stubAuth struct{ name string }

func (s stubAuth) NewToken(auth.UserClaims) (string, error)    { return "", nil }
func (s stubAuth) SetUserContext(echo.Context, string) error   { return nil }
func (s stubAuth) SetTokenCookies(echo.Context, string) string { return "" }
func (s stubAuth) ClearTokenCookies(echo.Context)              {}
func (s stubAuth) GetSigningKey() []byte                       { return nil }
func (s stubAuth) GetCookieDomain() string                     { return "" }
func (s stubAuth) GetTokenLifetime() time.Duration             { return 0 }
func (s stubAuth) GetMethod() auth.Method                      { return "" }
func (s stubAuth) GetUserFromContext(echo.Context) (*auth.User, error) {
	return &auth.User{UserClaims: auth.UserClaims{Name: s.name}}, nil
}

func newEnv(t *testing.T, as string) *env.Env {
	e := env.NewTestingEnvironment()
	if err := e.DB.AutoMigrate(&models.User{}, &models.AppAccess{}); err != nil {
		t.Fatal(err)
	}
	e.SetupCompleted = true
	e.Auth = stubAuth{name: as}
	for _, u := range []models.User{
		{Name: "root", Role: "admin"},
		{Name: "dev", Role: "developer"},
		{Name: "view", Role: "viewer"},
	} {
		e.DB.Create(&u)
	}
	e.DB.Create(&models.AppAccess{UserName: "dev", AppName: "alpha"})
	e.DB.Create(&models.AppAccess{UserName: "view", AppName: "alpha"})
	return e
}

func call(t *testing.T, as, method, target, body string) int {
	e := newEnv(t, as)
	ok := func(c echo.Context) error {
		// the handler must still be able to read the JSON body
		if body != "" {
			b := new(strings.Builder)
			buf := make([]byte, 512)
			n, _ := c.Request().Body.Read(buf)
			b.Write(buf[:n])
			if b.String() != body {
				t.Errorf("handler saw body %q, want %q", b.String(), body)
			}
		}
		return c.NoContent(http.StatusOK)
	}
	h := Authorize(e)(ok)
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	err := h(e.Router.NewContext(req, rec))
	if he, isHTTP := err.(*echo.HTTPError); isHTTP {
		return he.Code
	}
	return rec.Code
}

func TestAuthorize(t *testing.T) {
	cases := []struct {
		name               string
		as, method, target string
		body               string
		want               int
	}{
		{"admin anywhere", "root", "POST", "/api/apps/destroy", `{"name":"zzz"}`, 200},
		{"admin settings", "root", "GET", "/api/settings/users", "", 200},
		{"dev restarts assigned app", "dev", "POST", "/api/apps/restart", `{"name":"alpha"}`, 200},
		{"dev blocked from other app", "dev", "POST", "/api/apps/restart", `{"name":"beta"}`, 403},
		{"dev query-scoped read ok", "dev", "GET", "/api/apps/overview?name=alpha", "", 200},
		{"dev query-scoped other app", "dev", "GET", "/api/apps/overview?name=beta", "", 403},
		{"dev reads env of own app", "dev", "GET", "/api/apps/config?name=alpha", "", 200},
		{"dev cannot destroy", "dev", "POST", "/api/apps/destroy", `{"name":"alpha"}`, 403},
		{"dev cannot open settings", "dev", "GET", "/api/settings/users", "", 403},
		{"dev missing app name", "dev", "POST", "/api/apps/restart", `{}`, 403},
		{"viewer reads status", "view", "GET", "/api/apps/overview?name=alpha", "", 200},
		{"viewer cannot read env", "view", "GET", "/api/apps/config?name=alpha", "", 403},
		{"viewer cannot read logs", "view", "GET", "/api/apps/logs?name=alpha", "", 403},
		{"viewer cannot restart", "view", "POST", "/api/apps/restart", `{"name":"alpha"}`, 403},
		{"viewer lists apps", "view", "GET", "/api/apps/list", "", 200},
		{"unknown route needs admin", "dev", "GET", "/api/whatever", "", 403},
		{"deleted user locked out", "ghost", "GET", "/api/apps/list", "", 403},
	}
	for _, c := range cases {
		if got := call(t, c.as, c.method, c.target, c.body); got != c.want {
			t.Errorf("%s: %s %s as %s = %d, want %d", c.name, c.method, c.target, c.as, got, c.want)
		}
	}
}

func TestFilterApps(t *testing.T) {
	e := newEnv(t, "dev")
	req := httptest.NewRequest("GET", "/api/apps/list", nil)
	c := e.Router.NewContext(req, httptest.NewRecorder())
	loadAccess(e, c)
	got := FilterApps(c, []string{"alpha", "beta"})
	if len(got) != 1 || got[0] != "alpha" {
		t.Fatalf("dev sees %v, want [alpha]", got)
	}

	e = newEnv(t, "root")
	c = e.Router.NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())
	loadAccess(e, c)
	if got := FilterApps(c, []string{"alpha", "beta"}); len(got) != 2 {
		t.Fatalf("admin sees %v, want both", got)
	}
}
