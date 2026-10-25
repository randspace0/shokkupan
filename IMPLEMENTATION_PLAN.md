# Plan: Dokku Pro feature parity and an upgrade of the Svelte UI

Goal: shokkupan matches the feature set of [Dokku Pro](https://pro.dokku.com/), and its existing SvelteKit UI is upgraded to current versions of Svelte, SvelteKit, Tailwind CSS and daisyUI. There is no framework rewrite.

This file is a working plan. Update each stage's **Status** as work lands, and delete the file when every stage is complete.

## How this plan was made

- Feature list: the Dokku Pro marketing page plus the labels visible in its public demo (`pro.dokku.com/demo`), and the feature list from the maintainer below.
- Only text labels and structure were read from the demo, to learn what screens exist. Its code, styles and assets are proprietary and must not be copied. Everything here is built from scratch.
- "Today" below is from reading this repo and the open PRs. Anything marked *verify* is a belief that Stage 1 must confirm before work is estimated.

## Target feature set

### Apps
| Area | Detail | Today |
|---|---|---|
| Navbar actions | Visit, restart, rebuild, sync from git, stop, start | Exists in the Svelte `AppHeader`. *Verify* "sync from git" maps to `setup/sync-repo`. |
| Overview | Created at, updated at, git sha, deploy locked? | `GET /apps/info` has created_at and is_locked. Updated at and git sha: *verify*, likely missing. |
| Domains | Add, remove, enable and disable | Exists. |
| TLS | Let's Encrypt on or off per app | Exists (`/apps/letsencrypt`). Custom certs: not present. |
| Processes | Scale, resources, deploy checks | Exists. |
| Environment variables | View as rows or as a raw `.env` text | Rows exist. `.env` view and bulk edit are new UI. |
| Logs | Live tail, with filtering | Logs page exists. Filtering is new. |
| Settings | Rename, delete (double confirmation) | Exist. Double confirmation is new. |
| Settings | HTTP auth | Missing. Needs a Dokku plugin (see Stage 4). |
| Settings | GitHub webhook setup | Partly exists (`api/github.go`, `/github/events`). *Verify* the flow. |
| Settings | Linked services | Exists (`/apps/services`). |
| Also open | Notifications (#2), monitoring and charts (#4, #5), cron (#3) | New Svelte pages in open PRs. |

### Services
| Area | Detail | Today |
|---|---|---|
| Overview | Status, docker image, internal IP, data dir, config dir, exposed ports, service dir | *Verify* which fields the service API returns. |
| Logs | Live tail, as for apps | Exists. |
| Settings | Linked apps, link and unlink | Exists in part. |
| Create | Datastore type, image, version, config options, custom env | Exists in part. *Verify*. |

### Dokku Pro features not on the list above
These are in Dokku Pro. Decide in Stage 1 whether they are in scope; they are planned in Stage 4.
- Long-lived **API tokens** for scripting the JSON API (the API exists today, but only browser sessions can use it).
- **Users, teams and per-app access** (today: GitHub or password login, no roles).
- **Audit trail**: every command run on a user's behalf, with target, outcome, time and output.
- **Password reset links** and SSH key management per user (SSH keys exist in settings).
- **HTTP(S) git push** for deploying from behind a firewall.

## Architecture decisions

- **Keep SvelteKit** (SPA, `adapter-static`, embedded in the Go binary as today). No Node server in production, so the Docker image and `bootstrap.sh` stay as they are.
- **Tailwind CSS + daisyUI only.** No other component library.
- **Upgrade in small PRs**, each building and passing CI on its own: SvelteKit 1 to 2, Svelte 3 to 5, the query library, then Tailwind 3 to 4 with daisyUI 2 to 5. Use the official migration tools (`sv migrate`, the Tailwind upgrade tool) and read each release's migration guide first.
- **Known breaking changes to plan for**: daisyUI 2 to 5 drops the `btn loading` class (many `class:loading` uses), renames several components, and changes theming; `@sveltestack/svelte-query` is unmaintained and must move to `@tanstack/svelte-query`, whose API differs across roughly 60 files.
- **Auth stays in Go.** Login (GitHub OAuth, password and TOTP), JWT cookies and CSRF are implemented in `internal/server`. A JavaScript auth library such as Better Auth would need a Node runtime and its own database, which breaks the single-binary deployment. New auth features (API tokens, teams, audit) are added in Go.
- Backend additions follow the repo's conventions: handlers in `internal/server/api/<area>`, DTOs in `internal/server/dto`, tests next to the code, a bash plugin under `plugins/` when Dokku itself must do the work.

## Stages

### Stage 1: Upgrade the UI stack and audit
**Goal**: The Svelte UI runs on current dependencies with no behavior change, and the *verify* items are answered.
**Success Criteria**:
- PR per step: SvelteKit 2, Svelte 5, `@tanstack/svelte-query`, Tailwind 4 with daisyUI 5. Each builds with `vite build` and passes CI.
- Every existing page still works (login, setup, apps, services, settings), checked by a Playwright smoke test against the Go server with the mock Dokku client.
- The *verify* items above are resolved and the tables updated; decision recorded on which Pro extras (below) are in scope.
**Tests**: `vite build` and a Playwright smoke test (login, list apps, open an app) in CI; `go test ./...` unchanged.
**Status**: Not Started

### Stage 2: Apps parity
**Goal**: Every app screen from the table above is present in the Svelte UI.
**Success Criteria**:
- Navbar actions with live command output.
- Overview (including updated at and git sha, adding them to `GET /apps/info` if missing), domains, TLS, processes, environment (rows and `.env` views), logs with filtering, settings with double-confirmed delete, GitHub webhook, linked services.
- The Notifications, Cron and Monitoring pages from #2, #3, #4 and #5 (already Svelte; only restyle if daisyUI 5 breaks them).
**Tests**: Playwright flows per screen; Go handler tests for any new endpoint fields; log filter unit tests.
**Status**: Not Started

### Stage 3: Services parity
**Goal**: Every service screen from the table above works.
**Success Criteria**: service list and create flow (type, image, version, options, env); overview with all listed fields (extend the service API where it is short); logs; linked apps with link and unlink; existing backup pages ported.
**Tests**: Playwright flows; Go tests for new service-info fields.
**Status**: Not Started

### Stage 4: Dokku Pro extras
**Goal**: The features that make shokkupan comparable to Pro beyond the screens above. Each is its own PR and can be dropped without blocking Stage 5.
**Success Criteria**:
1. **API tokens**: table, `Authorization: Bearer` support, create and revoke UI, `curl` documentation.
2. **HTTP auth** per app: `plugins/http-auth` (Dokku has no core command for it), settings UI, install through `plugins/install.sh`.
3. **Users and teams**: admin and member roles, per-app and per-service grants, enforced in middleware on every `/api/apps` and `/api/services` route, with tests for denial.
4. **Audit trail**: record each command run through `commands.RequestExecution` with user, target, outcome, time and output; per-user view visible to that user and admins.
5. **HTTPS git push**: smart-HTTP endpoint authenticated with an API token, forwarding to Dokku. Highest risk; needs a real Dokku host to test.
**Tests**: Go unit tests for each; permission-denial tests for teams; plugin scripts get `bash -n` in CI and a manual host checklist in the PR.
**Status**: Not Started

### Stage 5: Polish and release
**Goal**: Parity work is documented and released.
**Success Criteria**: README and screenshots refreshed; API token and plugin docs published; release notes call out the dependency upgrades.
**Tests**: full Playwright suite green; image builds; manual upgrade from the previous release on a real host.
**Status**: Not Started

## Risks and open questions

- **Scope.** This is a full UI rewrite plus a permissions system. Stage 4 item 3 touches every route; ship it last and behind tests.
- **Open PRs.** #2 to #5 add Svelte pages. Merge them before the daisyUI upgrade, or they need fixing for it too.
- **Upgrade size.** The query-library swap and the daisyUI 5 class changes touch most components; keep them as separate PRs.
- **No live host.** Nothing here has been run against a real Dokku server. Plugin and Dokku-output assumptions need a staging host before release.
- **Licensing.** The upstream project is AGPL; keep `LICENSE` and attribution.
- **Questions for the maintainer**: which Stage 4 items are required for "on par"? Preferred test host for Stage 4 item 5?
