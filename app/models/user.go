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

// ExistsByEmail checks if a user exists with the given email
func (m *UserModel) ExistsByEmail(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	err := m.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %v", err)
	}

	return count > 0, nil
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
	// First, check if this is the main admin user (protect main admin)
	var email string
	err := m.DB.QueryRow("SELECT email FROM users WHERE id = ?", id).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user email: %v", err)
	}

	// Prevent deletion of main admin
	if email == "admin@example.com" {
		return fmt.Errorf("cannot delete the main administrator account")
	}

	// First, delete all blogs by this user
	_, err = m.DB.Exec("DELETE FROM blogs WHERE user_id = ?", id)
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

// Count returns the total number of users
func (m *UserModel) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`

	err := m.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %v", err)
	}

	return count, nil
}

// CountByRole returns the number of users with a specific role
func (m *UserModel) CountByRole(role string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE role = ?`

	err := m.DB.QueryRow(query, role).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %v", err)
	}

	return count, nil
}

// Update updates a user's information
func (m *UserModel) Update(id int, name, email, role string, password *string) error {
	// Check if email is unique (excluding current user)
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND id != ?`
	err := m.DB.QueryRow(query, email, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("email already exists")
	}

	// Update user information
	if password != nil && *password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		query = `UPDATE users SET name = ?, email = ?, role = ?, password = ?, updated_at = NOW() 
				 WHERE id = ?`
		_, err = m.DB.Exec(query, name, email, role, string(hashedPassword), id)
	} else {
		// Update without changing password
		query = `UPDATE users SET name = ?, email = ?, role = ?, updated_at = NOW() 
				 WHERE id = ?`
		_, err = m.DB.Exec(query, name, email, role, id)
	}

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// UpdateProfile updates the current user's profile with optional password verification
func (m *UserModel) UpdateProfile(userID int, name, email, currentPassword string, newPassword *string) error {
	// If changing password, verify current password
	if newPassword != nil && *newPassword != "" {
		if currentPassword == "" {
			return fmt.Errorf("current password is required when changing password")
		}

		user, err := m.GetByID(userID)
		if err != nil {
			return fmt.Errorf("user not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
		if err != nil {
			return fmt.Errorf("current password is incorrect")
		}
	}

	// Check if email is unique (excluding current user)
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND id != ?`
	err := m.DB.QueryRow(query, email, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("email already exists")
	}

	// Update user information
	if newPassword != nil && *newPassword != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		query = `UPDATE users SET name = ?, email = ?, password = ?, updated_at = NOW() 
				 WHERE id = ?`
		_, err = m.DB.Exec(query, name, email, string(hashedPassword), userID)
	} else {
		// Update without changing password
		query = `UPDATE users SET name = ?, email = ?, updated_at = NOW() 
				 WHERE id = ?`
		_, err = m.DB.Exec(query, name, email, userID)
	}

	if err != nil {
		return fmt.Errorf("failed to update profile: %v", err)
	}

	return nil
}
