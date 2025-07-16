// database/migrations/add_role_to_users.go - Add role column to users table
package main

import (
	"go-web-app/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	appConfig := config.LoadConfig()
	
	// Connect to database
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	log.Println("🔄 Adding role column to users table...")

	// Add role column to users table
	addRoleQuery := `
		ALTER TABLE users 
		ADD COLUMN role ENUM('admin', 'author', 'user') NOT NULL DEFAULT 'user'
	`

	_, err = db.Exec(addRoleQuery)
	if err != nil {
		log.Fatal("Failed to add role column: ", err)
	}

	log.Println("✅ Role column added to users table")

	// Update existing admin user to have admin role
	updateAdminQuery := `
		UPDATE users 
		SET role = 'admin' 
		WHERE email = 'admin@example.com'
	`

	_, err = db.Exec(updateAdminQuery)
	if err != nil {
		log.Printf("Warning: Failed to update admin user role: %v", err)
	} else {
		log.Println("✅ Admin user role updated")
	}

	// Update some users to author role
	updateAuthorQuery := `
		UPDATE users 
		SET role = 'author' 
		WHERE email IN ('john@example.com', 'jane@example.com')
	`

	_, err = db.Exec(updateAuthorQuery)
	if err != nil {
		log.Printf("Warning: Failed to update author users: %v", err)
	} else {
		log.Println("✅ Author users updated")
	}

	log.Println("🎉 Role migration completed successfully!")
}
