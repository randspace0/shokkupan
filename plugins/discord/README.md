# Dokku Discord Plugin

A lightweight Discord notification plugin for Dokku, written in pure bash. Get notified in your Discord channel whenever you deploy applications to your Dokku server.

## Features

- 🚀 Real-time deployment notifications to Discord
- 📝 Pure bash implementation - no dependencies required
- 🎯 Support for both global and per-app webhook configuration
- 👤 Includes commit author information when available
- ⚡ Lightweight and fast
- 🔧 Compatible with Dokku 0.19.0+ (tested with 0.37.5)

## Requirements

- Dokku 0.19.0 or later
- curl (usually pre-installed)

## Installation

```bash
# Install the plugin
curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s discord
```

The plugin will automatically install during the Dokku plugin installation process.

## Usage

### Setting up a webhook

First, create a Discord webhook:
1. Go to your Discord server settings
2. Navigate to Integrations → Webhooks
3. Click "New Webhook"
4. Configure your webhook and copy the webhook URL

### Configure globally (all apps)

```bash
dokku discord:set https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

### Configure per-app

```bash
dokku discord:set https://discord.com/api/webhooks/YOUR_WEBHOOK_URL myapp
```

### Get current webhook URL

```bash
# Get global webhook
dokku discord:get

# Get app-specific webhook
dokku discord:get myapp
```

### Clear webhook configuration

```bash
# Clear global webhook
dokku discord:clear

# Clear app-specific webhook
dokku discord:clear myapp
```

### Test webhook

Test your webhook configuration by sending a test notification:

```bash
# Test global webhook
dokku discord:test

# Test app-specific webhook
dokku discord:test myapp
```

This will send a test notification to your Discord channel to verify the webhook is configured correctly.

## How It Works

When you deploy an application, the plugin automatically sends two notifications to your Discord channel:

1. **Pre-deploy notification** (yellow/amber color) - Sent when deployment starts
2. **Post-deploy notification** (green color) - Sent when deployment completes

Each notification includes:

- 📱 Application name
- 🌐 Hostname
- 🔗 Application URL
- 👤 Commit author (if available)
- ⏰ Timestamp

Example pre-deploy notification:

```
John Doe <john@example.com>
🚀 Deployment starting...
Application: myapp
Hostname: example.com
URL: http://myapp.example.com
```

Example post-deploy notification:

```
John Doe <john@example.com>
Application: myapp
Hostname: example.com
URL: http://myapp.example.com
```

## Environment Variables

You can also enable Discord notifications per-app using environment variables:

```bash
dokku config:set myapp DISCORD_NOTIFY=1
```

This is useful when you want to use the global webhook but only notify for specific apps.

## Troubleshooting

### Notifications not being sent

1. Check if webhook is configured:
   ```bash
   dokku discord:get myapp
   ```

2. Verify the webhook URL is valid and the webhook still exists in Discord

3. Check Dokku logs:
   ```bash
   dokku logs myapp
   ```

4. Test the webhook manually:
   ```bash
   curl -X POST -H "Content-Type: application/json" \
     -d '{"embeds":[{"description":"Test","color":2407256}]}' \
     https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
   ```

### Permission issues

The plugin needs to write to the Dokku directory. Ensure the plugin has the correct permissions:

```bash
sudo chown -R dokku:dokku /var/lib/dokku/plugins/available/discord
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Credits

Based on the original [dokku-discord](https://github.com/llim5432/dokku-discord) plugin by llim5432, rewritten as a pure bash implementation for better compatibility and reduced dependencies.
