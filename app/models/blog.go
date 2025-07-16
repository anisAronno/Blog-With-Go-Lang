package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Blog represents a blog post in the system
type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"` // For displaying author name
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BlogModel handles blog database operations
type BlogModel struct {
	DB *sql.DB
}

// NewBlogModel creates a new BlogModel instance
func NewBlogModel(db *sql.DB) *BlogModel {
	return &BlogModel{DB: db}
}

// Create creates a new blog post in the database
func (m *BlogModel) Create(title, content string, userID int) (*Blog, error) {
	query := `INSERT INTO blogs (title, content, user_id, created_at, updated_at) 
			  VALUES (?, ?, ?, NOW(), NOW())`
	
	result, err := m.DB.Exec(query, title, content, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create blog: %v", err)
	}

	// Get the inserted blog ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get blog ID: %v", err)
	}

	// Return the created blog
	return m.GetByID(int(id))
}

// GetByID retrieves a blog by ID
func (m *BlogModel) GetByID(id int) (*Blog, error) {
	blog := &Blog{}
	query := `SELECT b.id, b.title, b.content, b.user_id, u.name as user_name,
			  b.created_at, b.updated_at 
			  FROM blogs b
			  LEFT JOIN users u ON b.user_id = u.id
			  WHERE b.id = ?`
	
	err := m.DB.QueryRow(query, id).Scan(
		&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.UserName,
		&blog.CreatedAt, &blog.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("blog not found")
		}
		return nil, fmt.Errorf("failed to get blog: %v", err)
	}

	return blog, nil
}

// GetAll retrieves all blog posts with pagination (public posts)
func (m *BlogModel) GetAll(limit, offset int) ([]*Blog, error) {
	query := `SELECT b.id, b.title, b.content, b.user_id, u.name as user_name,
			  b.created_at, b.updated_at 
			  FROM blogs b
			  LEFT JOIN users u ON b.user_id = u.id
			  ORDER BY b.created_at DESC
			  LIMIT ? OFFSET ?`
	
	rows, err := m.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get blogs: %v", err)
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		blog := &Blog{}
		err := rows.Scan(
			&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.UserName,
			&blog.CreatedAt, &blog.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blog: %v", err)
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetByUserID retrieves all blog posts by a specific user
func (m *BlogModel) GetByUserID(userID int) ([]*Blog, error) {
	query := `SELECT b.id, b.title, b.content, b.user_id, u.name as user_name,
			  b.created_at, b.updated_at 
			  FROM blogs b
			  LEFT JOIN users u ON b.user_id = u.id
			  WHERE b.user_id = ?
			  ORDER BY b.created_at DESC`
	
	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user blogs: %v", err)
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		blog := &Blog{}
		err := rows.Scan(
			&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.UserName,
			&blog.CreatedAt, &blog.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blog: %v", err)
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetAllBlogs retrieves all blog posts for admin users with pagination
func (m *BlogModel) GetAllBlogs(limit, offset int) ([]*Blog, error) {
	query := `SELECT b.id, b.title, b.content, b.user_id, u.name as user_name,
			  b.created_at, b.updated_at 
			  FROM blogs b
			  LEFT JOIN users u ON b.user_id = u.id
			  ORDER BY b.created_at DESC
			  LIMIT ? OFFSET ?`
	
	rows, err := m.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all blogs: %v", err)
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		blog := &Blog{}
		err := rows.Scan(
			&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.UserName,
			&blog.CreatedAt, &blog.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blog: %v", err)
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// Update updates a blog post
func (m *BlogModel) Update(id int, title, content string) (*Blog, error) {
	query := `UPDATE blogs SET title = ?, content = ?, updated_at = NOW() 
			  WHERE id = ?`
	
	_, err := m.DB.Exec(query, title, content, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog: %v", err)
	}

	// Return the updated blog
	return m.GetByID(id)
}

// Delete deletes a blog post
func (m *BlogModel) Delete(id int) error {
	query := `DELETE FROM blogs WHERE id = ?`
	
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("blog not found")
	}

	return nil
}

// Count returns the total number of blog posts
func (m *BlogModel) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM blogs`
	
	err := m.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count blogs: %v", err)
	}

	return count, nil
}

// CountUserBlogs returns the total number of blog posts for a specific user
func (m *BlogModel) CountUserBlogs(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM blogs WHERE user_id = ?`
	
	err := m.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user blogs: %v", err)
	}

	return count, nil
}

// CanUserEdit checks if a user can edit a specific blog post
func (m *BlogModel) CanUserEdit(blogID, userID int) (bool, error) {
	var ownerID int
	query := `SELECT user_id FROM blogs WHERE id = ?`
	
	err := m.DB.QueryRow(query, blogID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("blog not found")
		}
		return false, fmt.Errorf("failed to check blog ownership: %v", err)
	}

	return ownerID == userID, nil
}

// CanUserDelete checks if a user can delete a specific blog post (owner or admin)
func (m *BlogModel) CanUserDelete(blogID, userID int, userRole string) (bool, error) {
	// Admin users can delete any blog
	if userRole == "admin" {
		return true, nil
	}

	// Check if user owns the blog
	var ownerID int
	query := `SELECT user_id FROM blogs WHERE id = ?`
	
	err := m.DB.QueryRow(query, blogID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("blog not found")
		}
		return false, fmt.Errorf("failed to check blog ownership: %v", err)
	}

	return ownerID == userID, nil
}
