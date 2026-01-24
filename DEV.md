
# How to run the project in dev mode

- Start ONLY the database
`docker-compose -f docker-compose.dev.yml up -d`

- Run your Go app directly on your machine
`go run ./cmd/`

- Make changes to code, save, restart (Press Ctrl+C, then run again)
`go run ./cmd/`

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

## How to connect to your DB?
`psql -h localhost -p 5432 -U mygameshelf my_game_shelf`

> Replace with your project params

## How to run tests
 ```bash
go test ./src/api/handlers/...
```

Add `-v` to get more verbose output

`...` is go package wildcard, it means: “Run tests for this package and all subpackages recursively”


## How to run the web project

```bash
cd web
npm run dev
```

This will start the vite server