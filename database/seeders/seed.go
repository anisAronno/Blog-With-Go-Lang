// Package seeders handles database seeding
package seeders

import (
	"database/sql"
	"fmt"
	"log"
)

// Seeder interface for all seeders
type Seeder interface {
	Seed() error
	Clear() error
}

// SeederManager manages all seeders
type SeederManager struct {
	DB      *sql.DB
	seeders []Seeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(db *sql.DB) *SeederManager {
	manager := &SeederManager{
		DB: db,
	}
	manager.registerSeeders()
	return manager
}

// registerSeeders registers all available seeders
func (m *SeederManager) registerSeeders() {
	m.seeders = []Seeder{
		NewUserSeeder(m.DB),
		NewBlogSeeder(m.DB),
	}
}

// SeedAll runs all seeders
func (m *SeederManager) SeedAll() error {
	log.Println("🌱 Running all database seeders...")

	for _, seeder := range m.seeders {
		if err := seeder.Seed(); err != nil {
			return fmt.Errorf("seeder failed: %v", err)
		}
	}

	log.Println("✅ All seeders completed successfully")
	return nil
}

// ClearAll clears all seeded data
func (m *SeederManager) ClearAll() error {
	log.Println("🧹 Clearing all seeded data...")

	// Clear in reverse order to respect foreign key constraints
	for i := len(m.seeders) - 1; i >= 0; i-- {
		seeder := m.seeders[i]
		if err := seeder.Clear(); err != nil {
			return fmt.Errorf("clear failed: %v", err)
		}
	}

	log.Println("✅ All seeded data cleared")
	return nil
}

// RunSeeders is a convenience function to run all seeders
func RunSeeders(db *sql.DB) error {
	manager := NewSeederManager(db)
	return manager.SeedAll()
}
