// Package access decides which role may call which API route. It is pure
// (no database or HTTP types) so the whole policy is covered by a table test.
package access

import "strings"

const (
	RoleAdmin     = "admin"
	RoleDeveloper = "developer"
	RoleViewer    = "viewer"
)

// Roles lists the valid roles, most privileged first.
var Roles = []string{RoleAdmin, RoleDeveloper, RoleViewer}

func ValidRole(r string) bool {
	for _, v := range Roles {
		if r == v {
			return true
		}
	}
	return false
}

type Perm int

const (
	// PermRead: look at state (lists, status, domains, ...).
	PermRead Perm = iota
	// PermSecret: read data that can contain credentials (env values, logs,
	// build output).
	PermSecret
	// PermWrite: change an app (deploy, restart, env, domains, ...).
	PermWrite
	// PermAdmin: server-wide changes, users, creating or destroying things.
	PermAdmin
)

// Rule is what a route requires.
type Rule struct {
	Perm Perm
	// Scoped routes act on one app, which a non-admin must be assigned to.
	Scoped bool
}

// Allows reports whether role may use a route needing p. Unknown roles get
// nothing.
func Allows(role string, p Perm) bool {
	switch role {
	case RoleAdmin:
		return true
	case RoleDeveloper:
		return p <= PermWrite
	case RoleViewer:
		return p == PermRead
	}
	return false
}

// Classify maps a request to its Rule. Anything not listed needs PermAdmin,
// so a new route is locked down until someone classifies it.
func Classify(method, path string) Rule {
	path = strings.TrimSuffix(path, "/")
	get := method == "GET" || method == "HEAD"

	switch {
	case strings.HasPrefix(path, "/api/auth/"): // details, refresh, logout: every signed-in user
		return Rule{Perm: PermRead}

	case path == "/api/exec/status":
		return Rule{Perm: PermRead}

	case path == "/api/monitoring" || path == "/api/monitoring/server":
		return Rule{Perm: PermRead}
	case path == "/api/monitoring/stats":
		return Rule{Perm: PermRead, Scoped: true}

	case path == "/api/settings/versions" && get,
		path == "/api/settings/server-monitoring" && get:
		return Rule{Perm: PermRead}

	case path == "/api/notifications" || strings.HasPrefix(path, "/api/notifications/"):
		if get {
			return Rule{Perm: PermRead, Scoped: true}
		}
		return Rule{Perm: PermWrite, Scoped: true}

	case path == "/api/services/list", path == "/api/services/type",
		path == "/api/services/info", path == "/api/services/linked-apps":
		if get {
			return Rule{Perm: PermRead}
		}

	case strings.HasPrefix(path, "/api/apps/"):
		return classifyApps(get, strings.TrimPrefix(path, "/api/apps/"))
	}
	return Rule{Perm: PermAdmin}
}

func classifyApps(get bool, p string) Rule {
	switch p {
	case "list", "report":
		return Rule{Perm: PermRead} // filtered to assigned apps by the handler
	case "create", "destroy", "rename":
		return Rule{Perm: PermAdmin}
	case "config", "logs", "deploys/output":
		if get {
			return Rule{Perm: PermSecret, Scoped: true}
		}
	}
	if get {
		return Rule{Perm: PermRead, Scoped: true}
	}
	return Rule{Perm: PermWrite, Scoped: true}
}
