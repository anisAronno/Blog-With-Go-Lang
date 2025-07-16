// database/migrations/migrate.go - Database schema migration
package main

import (
	"go-web-app/config"
	"log"
)

func main() {
	// Load configuration and connect to database
	appConfig := config.LoadConfig()
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	log.Println("🔄 Running database migrations...")

	// Create users table
	userTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_email (email)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, err = db.Exec(userTableSQL)
	if err != nil {
		log.Fatal("Failed to create users table: ", err)
	}
	log.Println("✅ Users table created/verified")

	// Create blogs table
	blogTableSQL := `
	CREATE TABLE IF NOT EXISTS blogs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_user_id (user_id),
		INDEX idx_created_at (created_at),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, err = db.Exec(blogTableSQL)
	if err != nil {
		log.Fatal("Failed to create blogs table: ", err)
	}
	log.Println("✅ Blogs table created/verified")

	log.Println("🎉 Database migration completed successfully!")
}
