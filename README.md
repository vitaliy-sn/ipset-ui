# ipset-ui

**ipset-ui** is a web interface for managing IP sets (ipset) on a server. The project consists of a Go backend service and a Vue 3 frontend application.

## Features

- View, create, and delete ipset sets
- Add, remove, and search IP addresses and CIDRs in sets
- Import a list of CIDR or IP addresses from a file (for example, exported from [ip2location Visitor Blocker](https://www.ip2location.com/free/visitor-blocker))
- Import IP addresses by domain name
- Whois queries for IP addresses directly from the interface
- Backup and restore ipset sets via the web interface
- Convenient notifications and action confirmations

## Import from File

You can import a list of CIDR or IP addresses from a file (for example, exported from [ip2location Visitor Blocker](https://www.ip2location.com/free/visitor-blocker)). The web interface allows you to upload and bulk add entries to an ipset set.

## Backup and Restore

You can create backups of ipset sets and restore them when needed. Backups are stored in a dedicated directory and managed through the web interface.

## Quick Start

### Local

1. Install frontend dependencies:
   ```
   cd web
   npm install
   ```
2. Build the frontend:
   ```
   npm run build
   ```
3. Build and run the backend:
   ```
   cd ../api
   go build -o ipset-ui ./cmd/main.go
   ./ipset-ui
   ```

### Docker

```
docker compose up -d --build
```

## Environment Variables

- `APP_PORT` — backend port (default 8080) (optional)
- `APP_HOST` — backend listen address (default 0.0.0.0) (optional)
- `STATIC_DIR` — path to frontend static files (optional)
- `IPSET_BACKUP_DIR` — directory for ipset backups (required)
- `FRONTEND_URL` — used during development to proxy API requests to the frontend server (for example, when running `npm run dev`) (optional)
