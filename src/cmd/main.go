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
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
)

// Abstract run function to allow easier testing of main logic
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	//Create dbURL and run migrations
	err, dbURL := initializeDatabase()
	if err != nil {
		return err
	}

	dbPool, err := connectToDatabase(dbURL)
	if err != nil {
		return err
	}

	defer dbPool.Close()

	// Initialize repositories
	boardGameRepo := repository.NewBoardGameRepository(dbPool)
	if err := api.InitServer(boardGameRepo); err != nil {
		return err
	}

	return nil
}

func initializeDatabase() (error, string) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load environment variables
	dbUser := config.GetEnv("DB_USER", "mygameshelf")
	dbPassword := config.GetEnv("DB_PASSWORD", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbName := config.GetEnv("DB_NAME", "my_game_shelf")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	// Build database URL
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Run migrations
	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err, "failed"
	}

	return nil, dbURL
}

func connectToDatabase(dbURL string) (*pgxpool.Pool, error) {
	// Connect to database
	log.Println("Connecting to database...")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	// Verify connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		return nil, err
	}
	log.Println("Database connection established")

	return dbPool, nil
}

func runMigrations(dbURL string) error {
	log.Println("Running database migrations...")
	err := db.RunMigrations(dbURL)
	return err
}
