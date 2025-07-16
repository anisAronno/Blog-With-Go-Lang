package migrations

import (
	"database/sql"
	"fmt"
)

// CreateUsersTable creates the users table
func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role ENUM('admin', 'author', 'user') DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	fmt.Println("✅ Users table created successfully")
	return nil
}

// DropUsersTable drops the users table
func DropUsersTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS users;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop users table: %v", err)
	}

	fmt.Println("❌ Users table dropped successfully")
	return nil
}
