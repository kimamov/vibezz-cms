# Vibezz CMS

A modern, Kirby-inspired CMS built with Go, Gin, PostgreSQL, sqlc, and Nuxt.

## Architecture

- **Backend**: Go + Gin REST API
- **Database**: PostgreSQL with sqlc-generated typed queries
- **Admin Panel**: Nuxt 3 + Nuxt UI
- **Public Frontend**: Nuxt 3 SSR

## Prerequisites

- Go 1.21+
- Node.js 20+
- PostgreSQL 15+ (or Docker)
- [sqlc](https://sqlc.dev/) (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

## Quick Start

### 1. Start PostgreSQL

```bash
docker compose up -d
```

### 2. Copy environment config

```bash
cp .env.example .env
```

### 3. Run the API

```bash
cd apps/api
go run ./cmd/server
```

The API starts on `http://localhost:8080`. It automatically runs migrations and creates a default admin user:

- **Email**: admin@vibezz.cms
- **Password**: admin1234

### 4. Run the Admin Panel

```bash
npm run dev:admin
```

Opens on `http://localhost:3001`.

### 5. Run the Public Frontend

```bash
npm run dev:web
```

Opens on `http://localhost:3000`.

## Project Structure

```
vibezz-cms/
├── apps/
│   ├── api/              # Go/Gin backend
│   │   ├── cmd/server/   # Entry point
│   │   ├── internal/
│   │   │   ├── auth/     # Password hashing, JWT
│   │   │   ├── config/   # Environment configuration
│   │   │   ├── content/  # Content types, entries, users
│   │   │   ├── db/       # Database pool, migrations, sqlc generated code
│   │   │   ├── http/     # Router, handlers, middleware
│   │   │   └── media/    # File upload service
│   │   └── sqlc.yaml
│   ├── admin/            # Nuxt admin panel
│   └── web/              # Nuxt public frontend
├── db/
│   ├── migrations/       # SQL migration files
│   └── queries/          # sqlc query files
├── packages/
│   └── types/            # Shared TypeScript types
├── docker-compose.yml
└── Makefile
```

## API Endpoints

### Public

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/public/routes/*path` | Resolve content by URL path |
| GET | `/api/public/navigation` | Get site navigation tree |
| GET | `/api/public/media/:id` | Serve media file |

### Admin (requires auth)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/admin/auth/login` | Login |
| POST | `/api/admin/auth/refresh` | Refresh token |
| GET | `/api/admin/me` | Current user |
| GET/POST | `/api/admin/content-types` | List/Create content types |
| GET/PATCH/DELETE | `/api/admin/content-types/:id` | Read/Update/Delete content type |
| GET/POST | `/api/admin/entries` | List/Create entries |
| GET/PATCH/DELETE | `/api/admin/entries/:id` | Read/Update/Delete entry |
| POST | `/api/admin/entries/:id/publish` | Publish entry |
| POST | `/api/admin/entries/:id/unpublish` | Unpublish entry |
| GET/POST | `/api/admin/media` | List/Upload media |
| GET/DELETE | `/api/admin/media/:id` | Read/Delete media |

## Development

### Regenerate sqlc types

```bash
cd apps/api && sqlc generate
```

### Build the API binary

```bash
make build
```

### Run migrations only

```bash
make migrate
```

## Default Roles

- **admin**: Full access to all features
- **editor**: Can create, edit, publish/unpublish entries
- **author**: Can create and edit own entries
