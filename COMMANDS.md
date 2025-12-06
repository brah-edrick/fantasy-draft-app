# Useful Development Commands

Since we are running everything via Docker to keep your local machine clean, here are the commands you'll use often.

## 1. Lifecycle Management
**Start the App (and Database)**
```bash
docker-compose up --build
```
*   `--build`: Rebuilds the container (important if you changed `go.mod` or `Dockerfile`).
*   `-d`: Detached mode (runs in background).

**Stop Everything**
```bash
docker-compose down
```
*   Add `-v` (`docker-compose down -v`) to **delete the database volume** (start fresh).

**View Logs**
```bash
docker-compose logs -f api
```
*   `-f`: Follow (stream) the logs.

## 2. Dependencies (Go Modules)
We're assuming you want this to run in the container, not your local machine.

**Tidy Dependencies (Install new packages)**
Run this after adding an import to your code or changing `go.mod`:
```bash
docker run --rm -v "$PWD":/app -w /app golang:1.23-alpine go mod tidy
```
*   This mounts your folder, runs `go mod tidy`, and updates your `go.mod`/`go.sum` files locally.

## 3. Database & SQLC
**Generate Go Code from SQL**
Run this whenever you edit `db/query.sql` or `db/schema.sql`:
```bash
docker run --rm -v "$PWD":/src -w /src sqlc/sqlc generate
```

**Connect to DB via CLI**
```bash
docker exec -it fantasy-draft-app-db-1 psql -U fantasy_user -d fantasy_db
```
*   Password is: `secret_password`

## 4. Troubleshooting
**Rebuild everything from scratch**
If things get weird, nuke it and restart:
```bash
docker-compose down -v
docker-compose up --build
```
