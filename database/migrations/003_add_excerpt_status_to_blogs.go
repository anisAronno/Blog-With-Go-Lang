package migrations

import (
	"database/sql"
	"fmt"
)

// AddExcerptAndStatusToBlogs adds excerpt and status columns to blogs table
func AddExcerptAndStatusToBlogs(db *sql.DB) error {
	// Check if columns already exist
	var count int
	checkQuery := `SELECT COUNT(*) FROM information_schema.columns 
				  WHERE table_schema = DATABASE() AND table_name = 'blogs' AND column_name IN ('excerpt', 'status')`
	err := db.QueryRow(checkQuery).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing columns: %v", err)
	}

	if count > 0 {
		fmt.Println("⏭️  Excerpt and status columns already exist, skipping")
		return nil
	}

	query := `
	ALTER TABLE blogs 
	ADD COLUMN excerpt TEXT,
	ADD COLUMN status ENUM('draft', 'published') DEFAULT 'draft'
	`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to add excerpt and status columns to blogs table: %v", err)
	}

	fmt.Println("✅ Excerpt and status columns added to blogs table")
	return nil
}

// RemoveExcerptAndStatusFromBlogs removes excerpt and status columns from blogs table
func RemoveExcerptAndStatusFromBlogs(db *sql.DB) error {
	query := `ALTER TABLE blogs DROP COLUMN IF EXISTS excerpt, DROP COLUMN IF EXISTS status`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to remove excerpt and status columns from blogs table: %v", err)
	}

	fmt.Println("❌ Excerpt and status columns removed from blogs table")
	return nil
}
