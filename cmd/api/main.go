package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eddiarnoldo/my-game-shelf/db"
	"github.com/eddiarnoldo/my-game-shelf/internal/handlers"
	"github.com/eddiarnoldo/my-game-shelf/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	if err := startServer(boardGameRepo); err != nil {
		return err
	}

	return nil
}

func startServer(boardGameRepo *repository.BoardGameRepository) error {
	//Create gin router
	r := gin.Default()

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "*")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     parseOrigins(allowedOrigins),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	boardGameHandler := handlers.NewBoardGameHandler(boardGameRepo)

	setupRoutes(r, boardGameHandler)

	// Start server
	port := getEnv("APP_PORT", "8080")
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return err
	}
	return nil
}

func parseOrigins(origins string) []string {
	if origins == "*" {
		return []string{"*"}
	}
	parts := strings.Split(origins, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	return result
}

func initializeDatabase() (error, string) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load environment variables
	dbUser := getEnv("DB_USER", "mygameshelf")
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
