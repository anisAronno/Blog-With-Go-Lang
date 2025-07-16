// cmd/migrate/main.go - Migration CLI tool
package main

import (
	"flag"
	"fmt"
	"go-web-app/config"
	"go-web-app/database/migrations"
	"log"
	"os"
)

func main() {
	var action = flag.String("action", "up", "Migration action: up, down, status")
	flag.Parse()

	// Load configuration and connect to database
	appConfig := config.LoadConfig()
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Create migration manager
	manager := migrations.NewMigrationManager(db)

	// Execute based on action
	switch *action {
	case "up":
		if err := manager.Up(); err != nil {
			log.Fatal("Migration failed: ", err)
		}
	case "down":
		if err := manager.Down(); err != nil {
			log.Fatal("Migration rollback failed: ", err)
		}
	case "status":
		if err := manager.Status(); err != nil {
			log.Fatal("Failed to get migration status: ", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status")
		os.Exit(1)
	}
}
