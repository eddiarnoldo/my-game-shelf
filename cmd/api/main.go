package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Abstract run function to allow easier testing of main logic
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	//Create gin router
	r := gin.Default()
	setupRoutes(r)

	// Start server
	fmt.Println("Starting server on :8080")
	return r.Run() // listen and serve on
}
