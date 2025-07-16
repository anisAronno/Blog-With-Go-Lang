// cmd/seed/main.go - Seeder CLI tool
package main

import (
	"flag"
	"fmt"
	"go-web-app/config"
	"go-web-app/database/seeders"
	"log"
	"os"
)

func main() {
	var action = flag.String("action", "seed", "Seeder action: seed, clear")
	flag.Parse()

	// Load configuration and connect to database
	appConfig := config.LoadConfig()
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Create seeder manager
	manager := seeders.NewSeederManager(db)

	// Execute based on action
	switch *action {
	case "seed":
		if err := manager.SeedAll(); err != nil {
			log.Fatal("Seeding failed: ", err)
		}
	case "clear":
		if err := manager.ClearAll(); err != nil {
			log.Fatal("Clear failed: ", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: seed, clear")
		os.Exit(1)
	}
}
