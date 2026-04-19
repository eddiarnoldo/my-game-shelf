# My Game Shelf 🎲

A self-hosted board game collection tracker built with Go and React + TypeScript.

## Project Structure
```
my-game-shelf/
├── src/
│   ├── api/
│   │   ├── handlers/           # HTTP request handlers
│   │   │   ├── boardgame.go    # Board game CRUD endpoints
│   │   │   └── auth.go         # Auth endpoints (register/login/refresh/logout/invite)
│   │   ├── middleware/         # HTTP middleware
│   │   │   ├── cors.go         # CORS configuration
│   │   │   └── auth.go         # JWT validation, admin guard
│   │   ├── router/             # Route definitions
│   │   │   └── routes.go       # API route setup
│   │   └── api.go              # API initialization
│   │
│   ├── cmd/
│   │   └── main.go             # Application entry point
│   │
│   ├── config/
│   │   └── config.go           # Configuration management
│   │
│   ├── db/
│   │   ├── migrations/         # Database migrations (SQL)
│   │   │   ├── *_*.up.sql     # Migration scripts
│   │   │   └── *_*.down.sql   # Rollback scripts
│   │   └── migrate.go          # Migration runner
│   │
│   └── internal/
│       ├── models/             # Data models
│       │   ├── boardgame.go    # Board game model
│       │   └── auth.go         # User, Session, RefreshToken, Invite models
│       │
│       ├── repository/         # Database layer
│       │   ├── boardgame_repository.go
│       │   ├── user_repository.go
│       │   ├── session_repository.go
│       │   ├── refresh_token_repository.go
│       │   ├── invite_repository.go
│       │   └── errors.go
│       │
│       └── services/
│           └── email.go        # Email service (SMTP)
│
├── web/                        # React frontend (Vite + TypeScript)
│   ├── dist/                   # Production build output
│   ├── node_modules/           # Node dependencies
│   ├── public/                 # Static assets
│   ├── src/                    # React source code
│   │   ├── components/        # React components
│   │   ├── context/           # React context (auth state)
│   │   ├── pages/             # Page components
│   │   ├── App.tsx            # Main React component + routing
│   │   └── main.tsx           # React entry point
│   ├── .gitignore
│   ├── eslint.config.js       # ESLint configuration
│   ├── index.html             # HTML template
│   ├── package.json           # Node dependencies
│   ├── package-lock.json
│   └── README.md              # Vite/React README
│
├── .env                        # Environment variables (not in git)
├── .env.example                # Environment variables template
├── .gitignore                  # Git ignore rules
├── docker-compose.yml          # Production deployment
├── docker-compose.dev.yml      # Development (DB + Mailpit)
├── Dockerfile                  # Container build instructions
├── go.mod                      # Go module definition
├── go.sum                      # Go dependency checksums
└── README.md                   # This file
```

## Tech Stack

### Backend
- **Go 1.21+** - Backend API
- **Gin** - HTTP web framework
- **pgx** - PostgreSQL driver
- **golang-migrate** - Database migrations

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server

### Database
- **PostgreSQL 16** - Data storage

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 16 (via Docker)

### Development Setup

1. **Clone the repository**
```bash
   git clone https://github.com/eddiarnoldo/my-game-shelf.git
   cd my-game-shelf
```

2. **Configure environment**
```bash
   cp .env.example .env
   # Required: set DB_PASSWORD and JWT_SECRET
   # For email: set SMTP_* values (see Email Configuration below)
```

3. **Start the database and local mail server**
```bash
   docker compose -f docker-compose.dev.yml up -d
   # Starts PostgreSQL + Mailpit (local SMTP)
   # View captured emails at http://localhost:8025
```

4. **Run the backend**
```bash
   cd src
   go run ./cmd/main.go
   # API runs on http://localhost:8080
```

5. **Run the frontend**
```bash
   cd web
   npm install
   npm run dev
   # Frontend runs on http://localhost:5173
```

### API Endpoints

**Auth** (public)
- `POST /api/register` - Create account (requires invite code)
- `POST /api/login` - Login, returns access + refresh tokens
- `POST /api/refresh` - Rotate tokens
- `POST /api/logout` - Revoke session

**Admin only** (requires JWT + admin role)
- `POST /api/invites` - Generate and email an invite code

**Board games** (reads public, writes require JWT)
- `GET /api/boardgames` - List all board games
- `GET /api/boardgames/:id` - Get a specific board game
- `POST /api/boardgame` - Create a board game
- `PUT /api/boardgames/:id` - Update a board game
- `DELETE /api/boardgames/:id` - Delete a board game
- `POST /api/boardgame/:id/images` - Upload image
- `GET /api/boardgame/:id/images/coverThumbnail` - Get cover thumbnail
- `GET /api/boardgame/:id/image/:imageId` - Get image
- `GET /api/boardgame/:id/image/:imageId/thumbnail` - Get image thumbnail

## Folder Explanations

### `src/api/`
Contains the HTTP layer of the application:
- **handlers/** - Handle incoming HTTP requests and return responses
- **middleware/** - Process requests before they reach handlers (CORS, auth, logging)
- **router/** - Define which URLs map to which handlers

### `src/cmd/`
Contains the application entry point (`main.go`). This is where the program starts.

### `src/config/`
Configuration management - loads environment variables and application settings.

### `src/db/`
Database-related code:
- **migrations/** - SQL files that create/modify database schema
- **migrate.go** - Code that runs migrations automatically on startup

### `src/internal/`
Internal application code (not importable by external projects):
- **models/** - Define data structures (what a board game looks like)
- **repository/** - Database access layer (CRUD operations)

### `web/`
React frontend application built with Vite and TypeScript.

## Email Configuration

Invite codes are delivered by email. The setup differs between dev and production.

The active email backend is selected in `src/cmd/main.go` — two options are provided, one active, one commented out.

### Option 1 — Gmail (sends real emails)

Set in `src/.env`:
```
GMAIL_APP_PASSWORD=xxxx xxxx xxxx xxxx
```

Generate an App Password at **myaccount.google.com → Security → 2-Step Verification → App Passwords**. Use the 16-character code it provides. Do not use your regular Gmail password.

### Option 2 — Mailpit (local dev, no real sends)

`docker-compose.dev.yml` includes [Mailpit](https://github.com/axllent/mailpit) — a local SMTP server that captures all outgoing emails. No credentials needed. View captured emails at `http://localhost:8025`.

To switch: comment out the Gmail block in `src/cmd/main.go` and uncomment the Mailpit block. See `DEV.md` for the exact lines.

### Production (transactional email service)

For production with a custom domain, use Resend, Mailgun, or Brevo (all have free tiers). Use `services.NewSMTPEmailService()` with the provider's SMTP credentials.

## Self-Hosting

For production deployment:

1. **Configure environment**
```bash
   cp .env.example .env
   # Required: DB_PASSWORD, JWT_SECRET (32+ random chars), SMTP_* values
   # Set ALLOWED_ORIGINS to your domain or *
```

2. **Start with Docker Compose**
```bash
   docker-compose up -d
```

3. **Access the application API**
```
   http://YOUR_SERVER_IP:8080
```

## Development

### Create a new migration
```bash
migrate create -ext sql -dir src/db/migrations -seq your_migration_name
```

### Run migrations manually
```bash
migrate -path src/db/migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
  up
```

### Build frontend for production
```bash
cd web
npm run build
```

### Default account
On start there is a default account that will allow you to add new users.
- **username**: admin
- **password**: mygameshelf12345##@@

## License

MIT

## Contributing

Pull requests are welcome!