package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't include in JSON responses
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserModel handles user database operations
type UserModel struct {
	DB *sql.DB
}

// NewUserModel creates a new UserModel instance
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// Create creates a new user in the database
func (m *UserModel) Create(name, email, password string) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Insert user into database (default role is 'user')
	query := `INSERT INTO users (name, email, password, role, created_at, updated_at) 
			  VALUES (?, ?, ?, 'user', NOW(), NOW())`
	
	result, err := m.DB.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Get the inserted user ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	// Return the created user
	return m.GetByID(int(id))
}

// GetByID retrieves a user by ID
func (m *UserModel) GetByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password, role, created_at, updated_at 
			  FROM users WHERE id = ?`
	
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (m *UserModel) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password, role, created_at, updated_at 
			  FROM users WHERE email = ?`
	
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetAll retrieves all users
func (m *UserModel) GetAll() ([]*User, error) {
	query := `SELECT id, name, email, role, created_at, updated_at 
			  FROM users ORDER BY created_at DESC`
	
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Authenticate checks if the provided password matches the user's password
func (m *UserModel) Authenticate(email, password string) (*User, error) {
	user, err := m.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// EmailExists checks if an email already exists in the database
func (m *UserModel) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	
	err := m.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %v", err)
	}

	return count > 0, nil
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsAuthor checks if the user has author role
func (u *User) IsAuthor() bool {
	return u.Role == "author"
}

// CanManageBlogs checks if user can manage blogs (admin or author)
func (u *User) CanManageBlogs() bool {
	return u.Role == "admin" || u.Role == "author"
}

// Delete deletes a user by ID (admin only operation)
func (m *UserModel) Delete(id int) error {
	// First, delete all blogs by this user
	_, err := m.DB.Exec("DELETE FROM blogs WHERE user_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user blogs: %v", err)
	}

	// Then delete the user
	result, err := m.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
