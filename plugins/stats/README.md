# dokku-stats

Per-app CPU and memory usage for Dokku, written in pure bash. Built for the
[shokkupan](https://github.com/randspace0/shokkupan) monitoring page, but usable on its own.

## Requirements

- Dokku 0.19.0 or later
- The `dokku` user can run `docker` (true on any standard install)

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s stats
```

## Usage

```bash
# all apps
dokku stats:report

# one app
dokku stats:report myapp
```

One tab-separated line per running app container:

```
app  process-type  container-id  cpu%  used / limit
myapp  web  3f2a1b9c0d4e  0.35%  41.2MiB / 1.94GiB
```

Nothing is printed when no matching container is running. The plugin is
read-only and keeps no state; shokkupan samples it periodically.

## License

MIT License - see LICENSE file for details
