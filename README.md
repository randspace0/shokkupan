# shokkupan

![whaleboat logo](https://raw.githubusercontent.com/texm/shokku/main/web/static/images/whale-boat-sm.png)

_shokupan_ (食パン) is Japanese milk bread, the white loaf.

> A fork of [texm/shokku](https://github.com/texm/shokku), a friendly web interface for [Dokku](https://github.com/dokku/dokku), the open-source Platform-as-a-Service. Licensed AGPL, like the original.

## What it adds over shokku

- Rebranded and redesigned UI on Svelte 5, SvelteKit 2, Tailwind CSS 4 and daisyUI 5, with a signed-in user shown in the navbar.
- Raw `.env` editor next to the key/value table for app environment variables.
- Service pages skip Dokku service plugins that are not installed (e.g. mysql, mongo) instead of failing.
- **Monitoring** page: app status history and, with the `stats` plugin, per-app CPU and memory charts.
- **Server monitoring** (off by default): enable it under Settings to record host CPU, memory and network every 15s, chart them on the Monitoring page and log spikes (CPU or memory ≥ 90%, network ≥ 50 MB/s) as `sysmon: ... spike` in the container log. Network covers the container's own interfaces.
- Discord deploy notifications and alerts via the `discord` plugin.
- Cron UI.

## Install

Shokkupan runs as a Dokku app on the host it manages. See the original [installation docs](https://shokku.dev/docs) and `bootstrap.sh`. To deploy this fork to an existing Dokku app:

```bash
dokku builder:set shokku selected dockerfile
dokku git:sync --build shokku https://github.com/randspace0/shokkupan main
```

The first login goes through a GitHub App. Open `/setup` on your domain and enter the setup key printed in the container log (`dokku logs shokku | grep setup_key`). The GitHub App's callback URLs must include your host (`https://<host>/login/github/callback` and `https://<host>/setup/github/installed`), so recreate or edit the app if you change the domain.

Builds are heavy (npm install plus a Go compile). On a small host, build the image elsewhere and deploy it with `dokku git:from-image`.

SQLite access uses a single pooled connection so monitoring writes and sign-in
metadata reads cannot contend for locks within the server process. Keep one
server process per database file; separate processes do not share this pool.

## Plugins

Optional Dokku plugins live in [`plugins/`](plugins) and are installed on the Dokku host separately from shokkupan:

```bash
curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s <plugin>
```

| Plugin | Used for |
|---|---|
| `discord` | Deploy notifications and monitoring alerts |
| `stats` | CPU and memory charts on the Monitoring page |

## Development

```bash
cd ui && npm install && npm run dev   # UI on :5173, proxies /api to :5330
go build ./... && go test ./internal/...
```

`ui/` must be built (`npm run build`) before the Go binary can embed it. `ui/.mock/` is a gitignored place for a throwaway mock API.

## Support

Provided as-is, no support. For a supported product consider [Dokku Pro](https://pro.dokku.com).

## License

[AGPL](LICENSE). Original work © texm.
