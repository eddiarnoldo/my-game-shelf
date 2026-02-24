# My Game Shelf 🎲

A self-hosted board game collection tracker built with Go and React + TypeScript.

## Project Structure
```
my-game-shelf/
├── src/
│   ├── api/
│   │   ├── handlers/           # HTTP request handlers
│   │   │   └── boardgame.go    # Board game CRUD endpoints
│   │   ├── middleware/         # HTTP middleware
│   │   │   └── cors.go         # CORS configuration
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
│       │   └── boardgame.go   # Board game model
│       │
│       └── repository/         # Database layer
│           ├── boardgame_repository.go  # Board game data access
│           └── errors.go                # Repository errors
│
├── web/                        # React frontend (Vite + TypeScript)
│   ├── dist/                   # Production build output
│   ├── node_modules/           # Node dependencies
│   ├── public/                 # Static assets
│   ├── src/                    # React source code
│   │   ├── components/        # React components (if any)
│   │   ├── App.tsx            # Main React component
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
├── docker-compose.dev.yml      # Development database
├── Dockerfile                  # Container build instructions
├── go.mod                      # Go module definition
├── go.sum                      # Go dependency checksums
├── Makefile                    # Development commands (if exists)
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
   # Edit .env and set your DB_PASSWORD
```

3. **Start the database**
```bash
   docker compose -f docker-compose.dev.yml up -d
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

- `GET /api/boardgames` - List all board games
- `GET /api/boardgames/:id` - Get a specific board game
- `POST /api/boardgames` - Create a new board game
- `PUT /api/boardgames/:id` - Update a board game
- `DELETE /api/boardgames/:id` - Delete a board game

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

## Self-Hosting

For production deployment:

1. **Configure environment**
```bash
   cp .env.example .env
   # Set DB_PASSWORD and ALLOWED_ORIGINS=*
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

## License

MIT

## Contributing

Pull requests are welcome!