package seeders

import (
	"database/sql"
	"fmt"
	"go-web-app/app/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// UserSeeder handles user data seeding
type UserSeeder struct {
	DB        *sql.DB
	UserModel *models.UserModel
}

// NewUserSeeder creates a new user seeder
func NewUserSeeder(db *sql.DB) *UserSeeder {
	return &UserSeeder{
		DB:        db,
		UserModel: models.NewUserModel(db),
	}
}

// Seed creates default users
func (s *UserSeeder) Seed() error {
	log.Println("🌱 Seeding users...")

	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Admin User", "admin@example.com", "admin123", "admin"},
		{"John Author", "john@example.com", "password123", "author"},
		{"Jane User", "jane@example.com", "password123", "user"},
		{"Mike Writer", "mike@example.com", "password123", "author"},
	}

	for _, userData := range users {
		// Check if user already exists
		exists, err := s.UserModel.ExistsByEmail(userData.Email)
		if err != nil {
			return fmt.Errorf("failed to check if user exists: %v", err)
		}

		if exists {
			log.Printf("⏭️  User %s already exists, skipping", userData.Email)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		// Create user
		query := `INSERT INTO users (name, email, password, role, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, NOW(), NOW())`

		_, err = s.DB.Exec(query, userData.Name, userData.Email, string(hashedPassword), userData.Role)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", userData.Email, err)
		}

		log.Printf("✅ Created user: %s (%s)", userData.Name, userData.Role)
	}

	log.Println("✅ User seeding completed")
	return nil
}

// Clear removes all users (except admin for safety)
func (s *UserSeeder) Clear() error {
	log.Println("🧹 Clearing user data...")

	query := `DELETE FROM users WHERE role != 'admin'`
	_, err := s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear users: %v", err)
	}

	log.Println("✅ User data cleared (admin preserved)")
	return nil
}
