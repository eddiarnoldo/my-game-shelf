package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eddiarnoldo/my-game-shelf/src/api"
	"github.com/eddiarnoldo/my-game-shelf/src/config"
	"github.com/eddiarnoldo/my-game-shelf/src/db"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	err, dbURL := initializeDatabase()
	if err != nil {
		return err
	}

	dbPool, err := connectToDatabase(dbURL)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	boardGameRepo := repository.NewBoardGameRepository(dbPool)
	imageRepo := repository.NewBoardGameImageRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)
	sessionRepo := repository.NewSessionRepository(dbPool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(dbPool)
	inviteRepo := repository.NewInviteRepository(dbPool)

	// Option 1: Gmail — sends real emails to your inbox.
	// Requires GMAIL_APP_PASSWORD in .env (generate at myaccount.google.com → Security → App Passwords).
	gmailPass := config.GetEnv("GMAIL_APP_PASSWORD", "")
	log.Printf("DEBUG GMAIL_APP_PASSWORD=%q", gmailPass)
	emailService := services.NewGmailEmailService(
		"eddiarnoldo@gmail.com",
		gmailPass,
		config.GetEnv("APP_BASE_URL", "http://localhost:5173"),
	)

	// Option 2: Mailpit — captures emails locally, no real sends.
	// View at http://localhost:8025. Requires docker-compose.dev.yml to be running.
	// emailService := services.NewMailpitEmailService(
	// 	config.GetEnv("APP_BASE_URL", "http://localhost:5173"),
	// )

	if err := api.InitServer(
		boardGameRepo,
		imageRepo,
		userRepo,
		sessionRepo,
		refreshTokenRepo,
		inviteRepo,
		emailService,
	); err != nil {
		return err
	}

	return nil
}

func initializeDatabase() (error, string) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbUser := config.GetEnv("DB_USER", "mygameshelf")
	dbPassword := config.GetEnv("DB_PASSWORD", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbName := config.GetEnv("DB_NAME", "my_game_shelf")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err, "failed"
	}

	return nil, dbURL
}

func connectToDatabase(dbURL string) (*pgxpool.Pool, error) {
	log.Println("Connecting to database...")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		return nil, err
	}
	log.Println("Database connection established")

	return dbPool, nil
}

func runMigrations(dbURL string) error {
	log.Println("Running database migrations...")
	return db.RunMigrations(dbURL)
}
