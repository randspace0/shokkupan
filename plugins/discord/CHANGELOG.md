# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2026-01-12

### Changed
- Complete rewrite in pure bash (removed Node.js/TypeScript dependency)
- Simplified installation process - no build step required
- Updated to support Dokku 0.37.5
- Changed from `commands` wrapper to proper `subcommands/` structure

### Added
- Native bash implementation for all functionality
- Better integration with Dokku's plugin system
- Improved error handling and logging

### Removed
- Node.js/Bun dependency
- TypeScript build process
- npm/package.json dependencies

## [1.0.0] - 2024

### Added
- Initial TypeScript implementation
- Discord webhook notifications
- Global and per-app configuration
- Post-deploy hook integration
- Commit author information in notifications
