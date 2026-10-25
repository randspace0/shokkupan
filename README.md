# shokku
![whaleboat logo](https://raw.githubusercontent.com/texm/shokku/main/web/static/images/whale-boat-sm.png)

_(show-dokku)_

A friendly interface for [Dokku](https://github.com/dokku/dokku), the open-source Platform-as-a-Service.

![](https://raw.githubusercontent.com/texm/shokku/main/web/static/images/app-overview.webp)

## Usage
Read the [docs](https://shokku.dev/docs)

## Support
Shokku is provided 'as-is' - no support is available.
Consider [Dokku Pro](https://pro.dokku.com).

## License
[AGPL License](https://github.com/texm/shokku/blob/master/LICENSE)

## Plugins
Optional Dokku plugins live in [`plugins/`](plugins) and are installed on the Dokku host separately from shokku:

```bash
curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s <plugin>
```

| Plugin | Used for |
|---|---|
| `discord` | Deploy notifications and monitoring alerts |
