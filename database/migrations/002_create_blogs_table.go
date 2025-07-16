package migrations

import (
	"database/sql"
	"fmt"
)

// CreateBlogsTable creates the blogs table
func CreateBlogsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS blogs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create blogs table: %v", err)
	}

	fmt.Println("✅ Blogs table created successfully")
	return nil
}

// DropBlogsTable drops the blogs table
func DropBlogsTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS blogs;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop blogs table: %v", err)
	}

	fmt.Println("❌ Blogs table dropped successfully")
	return nil
}
