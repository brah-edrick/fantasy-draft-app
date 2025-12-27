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
docker run --rm -v "$PWD/server":/app -w /app golang:1.23-alpine go mod tidy
```
*   This mounts your server folder, runs `go mod tidy`, and updates your `go.mod`/`go.sum` files locally.

## 3. Database
**Connect to DB via CLI**
```bash
docker exec -it fantasy-draft-app-db-1 psql -U fantasy_user -d fantasy_db
```
*   Password is: `secret_password`

## 4. Database Seeding

**Step 1: Start the Database**
First, make sure the database is running:
```bash
docker-compose up -d db
```

**Step 2: Seed the Database**
Run the seed script from the synthetic-data folder:
```bash
cd server/synthetic-data
REAL_DATA_FILE="/Users/brandon/Projects/fantasy-draft-app/server/synthetic-data/real-data.json" \
DATABASE_URL="postgres://fantasy_user:secret_password@localhost:5432/fantasy_db?sslmode=disable" \
go run . seed
```

Or using Docker (if you don't want Go installed locally):
```bash
docker-compose exec api sh -c "REAL_DATA_FILE=/app/synthetic-data/real-data.json cd /app/synthetic-data && go run . seed"
```

**Step 3: Verify the Data**
Connect to the database and check the tables:
```bash
docker exec -it fantasy-draft-app-db-1 psql -U fantasy_user -d fantasy_db
```

Then run some queries:
```sql
-- Check table counts
SELECT 'conferences' as table_name, COUNT(*) FROM conferences
UNION ALL
SELECT 'divisions', COUNT(*) FROM divisions
UNION ALL
SELECT 'pro_teams', COUNT(*) FROM pro_teams
UNION ALL
SELECT 'players', COUNT(*) FROM players
UNION ALL
SELECT 'yearly_stats', COUNT(*) FROM yearly_stats;

-- View sample player data
SELECT first_name, last_name, position, skill FROM players LIMIT 10;

-- View sample yearly stats
SELECT p.first_name, p.last_name, ys.year, ys.stats
FROM yearly_stats ys
JOIN players p ON ys.player_id = p.id
LIMIT 5;
```

**Reset and Re-seed**
To completely reset the database and seed fresh:
```bash
# Stop and remove the database volume
docker-compose down -v

# Restart (this re-runs schema.sql)
docker-compose up -d db

# Wait for DB to be ready, then seed
sleep 5
cd server/synthetic-data
REAL_DATA_FILE="/Users/brandon/Projects/fantasy-draft-app/server/synthetic-data/real-data.json" \
DATABASE_URL="postgres://fantasy_user:secret_password@localhost:5432/fantasy_db?sslmode=disable" \
go run . seed
```

## 5. GraphQL

**Regenerate GraphQL Code**
Run this whenever you edit `server/graph/schema.graphql`:
```bash
cd server
go run github.com/99designs/gqlgen generate
```

**Start the GraphQL Server (locally)**
```bash
cd server
DATABASE_URL="postgres://fantasy_user:secret_password@localhost:5432/fantasy_db?sslmode=disable" \
go run ./cmd/server
```

**Start with Docker**
```bash
docker-compose up
```

Then visit:
- **GraphQL Playground**: http://localhost:8080/playground
- **GraphQL Endpoint**: http://localhost:8080/graphql

**Example Query**
```graphql
query {
  teams {
    city
    name
    players {
      firstName
      lastName
      position
    }
  }
}
```

## 6. Troubleshooting
**Rebuild everything from scratch**
If things get weird, nuke it and restart:
```bash
docker-compose down -v
docker-compose up --build
```
