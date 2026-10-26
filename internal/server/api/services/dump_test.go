package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func ctx() (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	return echo.New().NewContext(httptest.NewRequest("GET", "/", nil), rec), rec
}

func TestStreamDumpSendsData(t *testing.T) {
	c, rec := ctx()
	err := streamDump(c, strings.NewReader("PGDMP-data"), strings.NewReader(""), func() error { return nil }, "db.dump")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Body.String() != "PGDMP-data" {
		t.Fatalf("body = %q", rec.Body.String())
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, `filename="db.dump"`) {
		t.Fatalf("content-disposition = %q", cd)
	}
}

func TestStreamDumpEarlyFailureIsAnError(t *testing.T) {
	c, rec := ctx()
	err := streamDump(c, strings.NewReader(""), strings.NewReader("service does not exist"), func() error { return nil }, "x")
	var he *echo.HTTPError
	if !errors.As(err, &he) || he.Code != http.StatusBadGateway {
		t.Fatalf("want 502, got %v", err)
	}
	if rec.Body.Len() != 0 || rec.Header().Get("Content-Disposition") != "" {
		t.Fatal("must not start a download when export produced nothing")
	}
}
