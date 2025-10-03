# Loan Origination System (Go)
LOS service with layered architecture, chi router, and placeholders for repos, jobs, and notifications.

## Prerequisites
- Go 1.21+

## Run
bash
go run ./cmd/api


Server listens on port 8080 by default. Set a custom port with PORT:
bash
PORT=9090 go run ./cmd/api


## Endpoints (scaffold)
- GET /health
- POST /api/v1/loans
- PUT /api/v1/agents/{agent_id}/loans/{loan_id}/decision
- GET /api/v1/loans/status-count
- GET /api/v1/customers/top
- GET /api/v1/loans?status&size&page

## Configuration
- PORT (default: 8080)
- DATABASE_URL (default: postgres://postgres:postgres@localhost:5432/los?sslmode=disable)
- WORKER_COUNT (default: 4)

## Database & Migrations
- Ensure PostgreSQL is running and database los exists.
- Enable pgcrypto extension (migration includes CREATE EXTENSION IF NOT EXISTS pgcrypto;).
- Apply migrations in migrations/ (use your preferred tool, e.g., golang-migrate):
bash
migrate -path migrations -database "$DATABASE_URL" up


### Quick local Postgres (docker)
bash
docker run --name los-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=los -p 5432:5432 -d postgres:15
