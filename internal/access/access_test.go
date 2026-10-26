package access

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct {
		method, path string
		want         Rule
	}{
		{"GET", "/api/apps/list", Rule{PermRead, false}},
		{"GET", "/api/apps/report", Rule{PermRead, false}},
		{"GET", "/api/apps/overview", Rule{PermRead, true}},
		{"GET", "/api/apps/domains", Rule{PermRead, true}},
		{"GET", "/api/apps/config", Rule{PermSecret, true}},
		{"POST", "/api/apps/config", Rule{PermWrite, true}},
		{"GET", "/api/apps/logs", Rule{PermSecret, true}},
		{"GET", "/api/apps/deploys/output", Rule{PermSecret, true}},
		{"GET", "/api/apps/deploys", Rule{PermRead, true}},
		{"POST", "/api/apps/restart", Rule{PermWrite, true}},
		{"POST", "/api/apps/process/scale", Rule{PermWrite, true}},
		{"POST", "/api/apps/create", Rule{PermAdmin, false}},
		{"POST", "/api/apps/destroy", Rule{PermAdmin, false}},
		{"POST", "/api/apps/rename", Rule{PermAdmin, false}},
		{"GET", "/api/services/list", Rule{PermRead, false}},
		{"POST", "/api/services/destroy", Rule{PermAdmin, false}},
		{"GET", "/api/services/logs", Rule{PermAdmin, false}},
		{"GET", "/api/services/backups/report", Rule{PermAdmin, false}},
		{"GET", "/api/settings/users", Rule{PermAdmin, false}},
		{"GET", "/api/settings/versions", Rule{PermRead, false}},
		{"POST", "/api/settings/server-monitoring", Rule{PermAdmin, false}},
		{"GET", "/api/settings/ssh-keys", Rule{PermAdmin, false}},
		{"GET", "/api/notifications", Rule{PermRead, true}},
		{"POST", "/api/notifications", Rule{PermWrite, true}},
		{"GET", "/api/monitoring", Rule{PermRead, false}},
		{"GET", "/api/monitoring/stats", Rule{PermRead, true}},
		{"GET", "/api/exec/status", Rule{PermRead, false}},
		{"GET", "/api/auth/details", Rule{PermRead, false}},
		{"GET", "/api/brand-new-route", Rule{PermAdmin, false}}, // deny by default
		{"GET", "/api/apps/domains/", Rule{PermRead, true}},     // trailing slash
	}
	for _, c := range cases {
		if got := Classify(c.method, c.path); got != c.want {
			t.Errorf("%s %s = %+v, want %+v", c.method, c.path, got, c.want)
		}
	}
}

func TestAllows(t *testing.T) {
	perms := []Perm{PermRead, PermSecret, PermWrite, PermAdmin}
	want := map[string][]bool{
		RoleAdmin:     {true, true, true, true},
		RoleDeveloper: {true, true, true, false},
		RoleViewer:    {true, false, false, false},
		"":            {false, false, false, false},
		"root":        {false, false, false, false},
	}
	for role, exp := range want {
		for i, p := range perms {
			if got := Allows(role, p); got != exp[i] {
				t.Errorf("Allows(%q, %v) = %v, want %v", role, p, got, exp[i])
			}
		}
	}
}
