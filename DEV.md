
# How to run the project in dev mode

- Start ONLY the database
`docker-compose -f docker-compose.dev.yml up -d`

- Run your Go app directly on your machine
`go run ./cmd/api`

- Make changes to code, save, restart (Press Ctrl+C, then run again)
`go run ./cmd/api`

- When done, stop database
`docker-compose -f docker-compose.dev.yml down`

# DB Migrations

This project uses `golang-migrage`

```bash
Download for macOS
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.darwin-amd64.tar.gz | tar xvz

# Move to a location in your PATH
sudo mv migrate /usr/local/bin/

# Verify
migrate -version
```

## Go dependencies
```bash
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

## How to create a new migration?
Golang migrate creates up and down files so you can revert any migrations if needed

`migrate create -ext sql -dir db/migrations -seq create_games_table`

This will create 
```bash
db/migrations/000001_create_games_table.up.sql
db/migrations/000001_create_games_table.down.sql
```