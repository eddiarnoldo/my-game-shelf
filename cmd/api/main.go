package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eddiarnoldo/my-game-shelf/db"
	"github.com/eddiarnoldo/my-game-shelf/internal/handlers"
	"github.com/eddiarnoldo/my-game-shelf/internal/repository"
	"github.com/gin-gonic/gin"
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

	// Connect to database
	log.Println("Connecting to database...")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Verify connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Database connection established")

	// Initialize repositories
	boardGameRepo := repository.NewBoardGameRepository(dbPool)

	// Initialize handlers
	boardGameHandler := handlers.NewBoardGameHandler(boardGameRepo)

	//Create gin router
	r := gin.Default()
	setupRoutes(r, boardGameHandler)

	// Start server
	fmt.Println("Starting server on :8080")
	return r.Run() // listen and serve on
}

func initializeDatabase() (error, string) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load environment variables
	dbUser := getEnv("DB_USER", "gameshelf")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "my_game_shelf")

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

func runMigrations(dbURL string) error {
	log.Println("Running database migrations...")
	err := db.RunMigrations(dbURL)
	return err
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
