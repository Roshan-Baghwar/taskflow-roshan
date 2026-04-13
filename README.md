# TaskFlow — Backend (Zomato Assignment)

Minimal task management system

## Tech Stack
- **Backend**: Go 1.23 + Gin
- **Database**: PostgreSQL 16
- **Auth**: bcrypt + JWT (24h expiry)
- **Migrations**: golang-migrate (runs automatically on startup)
- **Logging**: slog
- **Docker**: Multi-stage build

## Architecture Decisions
- **Layered architecture** (handler → service → repository) for clean separation.
- Migrations run automatically on container start (simple & reliable for this scope).
- Pagination + `/stats` endpoint added as bonus.
- Structured JSON error responses exactly as specified.
- Owner-only checks for delete/update where required.
- Tradeoff: No full integration tests in this version (could add with `testify`).

## Running Locally
```bash
git clone https://github.com/Roshan-Baghwar/taskflow-roshan.git
cd taskflow-roshan
cp .env.example .env
docker compose up --build
```

App runs at http://localhost:8080

## Running Migrations
Migrations run automatically on backend startup. No manual step needed.

## Test Credentials
Email: test@example.com
Password: password123