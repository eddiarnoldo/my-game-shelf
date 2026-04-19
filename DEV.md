
# How to run the project in dev mode

- Start the database and the local mail server (Mailpit)
`docker compose -f docker-compose.dev.yml up -d`

- Run your Go app directly on your machine
`go run ./cmd/`

- Make changes to code, save, restart (Press Ctrl+C, then run again)
`go run ./cmd/`

- When done, stop services
`docker-compose -f docker-compose.dev.yml down`

## Email — switching between Gmail and Mailpit

`src/cmd/main.go` has two options. Only one should be active at a time.

**Option 1 — Gmail (sends real emails):**
```go
emailService := services.NewGmailEmailService(
    "eddiarnoldo@gmail.com",
    config.GetEnv("GMAIL_APP_PASSWORD", ""),
    config.GetEnv("APP_BASE_URL", "http://localhost:5173"),
)
```
Requires a Gmail App Password in `src/.env`:
```
GMAIL_APP_PASSWORD=xxxx xxxx xxxx xxxx
```
Generate one at: **myaccount.google.com → Security → 2-Step Verification → App Passwords**
(Name it anything, e.g. "my-game-shelf". Use the 16-char code it gives you.)

**Option 2 — Mailpit (local dev, no real sends):**
```go
// emailService := services.NewMailpitEmailService(
//     config.GetEnv("APP_BASE_URL", "http://localhost:5173"),
// )
```
Uncomment this block and comment out the Gmail block. Emails are captured at `http://localhost:8025`. No credentials needed.

## JWT Secret

The app requires a `JWT_SECRET` to sign authentication tokens. Generate one before first run:

```bash
openssl rand -base64 32
```

Paste the output into `src/.env`:

```
JWT_SECRET=your-generated-value-here
```

**Never commit this value to git.** If it leaks, all existing sessions become invalid after you rotate it (which is the correct response).

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