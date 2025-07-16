// Package migrations handles database schema migrations
package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	ID       string
	Name     string
	UpFunc   func(*sql.DB) error
	DownFunc func(*sql.DB) error
}

// MigrationManager manages database migrations
type MigrationManager struct {
	DB         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	manager := &MigrationManager{
		DB: db,
	}
	manager.registerMigrations()
	return manager
}

// registerMigrations registers all available migrations
func (m *MigrationManager) registerMigrations() {
	m.migrations = []Migration{
		{
			ID:       "001",
			Name:     "create_users_table",
			UpFunc:   CreateUsersTable,
			DownFunc: DropUsersTable,
		},
		{
			ID:       "002",
			Name:     "create_blogs_table",
			UpFunc:   CreateBlogsTable,
			DownFunc: DropBlogsTable,
		},
	}
}

// createMigrationsTable creates the migrations tracking table
func (m *MigrationManager) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, err := m.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}
	return nil
}

// isExecuted checks if a migration has been executed
func (m *MigrationManager) isExecuted(migrationID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM migrations WHERE id = ?`
	err := m.DB.QueryRow(query, migrationID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// markAsExecuted marks a migration as executed
func (m *MigrationManager) markAsExecuted(migration Migration) error {
	query := `INSERT INTO migrations (id, name) VALUES (?, ?)`
	_, err := m.DB.Exec(query, migration.ID, migration.Name)
	return err
}

// markAsReverted removes a migration from executed list
func (m *MigrationManager) markAsReverted(migrationID string) error {
	query := `DELETE FROM migrations WHERE id = ?`
	_, err := m.DB.Exec(query, migrationID)
	return err
}

// Up runs all pending migrations
func (m *MigrationManager) Up() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	fmt.Println("🚀 Running database migrations...")

	for _, migration := range m.migrations {
		executed, err := m.isExecuted(migration.ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if executed {
			fmt.Printf("⏭️  Skipping migration %s_%s (already executed)\n", migration.ID, migration.Name)
			continue
		}

		fmt.Printf("⬆️  Running migration %s_%s...\n", migration.ID, migration.Name)
		if err := migration.UpFunc(m.DB); err != nil {
			return fmt.Errorf("migration %s failed: %v", migration.ID, err)
		}

		if err := m.markAsExecuted(migration); err != nil {
			return fmt.Errorf("failed to mark migration as executed: %v", err)
		}
	}

	fmt.Println("✅ All migrations completed successfully")
	return nil
}

// Down reverts the last migration
func (m *MigrationManager) Down() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	// Find the last executed migration
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		executed, err := m.isExecuted(migration.ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if executed {
			fmt.Printf("⬇️  Reverting migration %s_%s...\n", migration.ID, migration.Name)
			if err := migration.DownFunc(m.DB); err != nil {
				return fmt.Errorf("failed to revert migration %s: %v", migration.ID, err)
			}

			if err := m.markAsReverted(migration.ID); err != nil {
				return fmt.Errorf("failed to mark migration as reverted: %v", err)
			}

			fmt.Println("✅ Migration reverted successfully")
			return nil
		}
	}

	fmt.Println("ℹ️  No migrations to revert")
	return nil
}

// Status shows the status of all migrations
func (m *MigrationManager) Status() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	fmt.Println("📊 Migration Status:")
	fmt.Println("-------------------")

	for _, migration := range m.migrations {
		executed, err := m.isExecuted(migration.ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		status := "❌ Pending"
		if executed {
			status = "✅ Executed"
		}

		fmt.Printf("%s %s_%s\n", status, migration.ID, migration.Name)
	}

	return nil
}

// RunMigrations is a convenience function to run all migrations
func RunMigrations(db *sql.DB) error {
	manager := NewMigrationManager(db)
	return manager.Up()
}

// Legacy function for backward compatibility
func CreateTables(db *sql.DB) error {
	log.Println("⚠️  CreateTables is deprecated. Use RunMigrations instead.")
	return RunMigrations(db)
}
