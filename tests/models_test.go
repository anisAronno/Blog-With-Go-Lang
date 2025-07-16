// tests/models_test.go - Unit tests for model functionality
package tests

import (
	"database/sql"
	"go-web-app/app/models"
	"go-web-app/config"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

// setupTestDB initializes a test database connection
func setupTestDB() {
	// Load test configuration
	testConfig := &config.Config{
		DBHost:     "localhost",
		DBPort:     "3308",
		DBName:     "go_web_app_test", // Use separate test database
		DBUser:     "root",
		DBPassword: "bs@123",
	}

	var err error
	testDB, err = config.ConnectDatabase(testConfig)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create test tables
	createTestTables()
}

// createTestTables creates necessary tables for testing
func createTestTables() {
	// Create users table
	userTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := testDB.Exec(userTableSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create blogs table
	blogTableSQL := `
	CREATE TABLE IF NOT EXISTS blogs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err = testDB.Exec(blogTableSQL)
	if err != nil {
		log.Fatalf("Failed to create blogs table: %v", err)
	}
}

// cleanupTestDB cleans up test data
func cleanupTestDB() {
	if testDB != nil {
		// Clean up test data
		testDB.Exec("DELETE FROM blogs")
		testDB.Exec("DELETE FROM users")
		testDB.Close()
	}
}

// TestUserModel tests user model functionality
func TestUserModel(t *testing.T) {
	setupTestDB()
	defer cleanupTestDB()

	userModel := models.NewUserModel(testDB)

	// Test user creation
	t.Run("CreateUser", func(t *testing.T) {
		user, err := userModel.Create("Test User", "test@example.com", "password123")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if user.Name != "Test User" {
			t.Errorf("Expected name 'Test User', got '%s'", user.Name)
		}

		if user.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
		}

		if user.ID == 0 {
			t.Error("Expected user ID to be set")
		}
	})

	// Test user authentication
	t.Run("AuthenticateUser", func(t *testing.T) {
		// Create a user first
		userModel.Create("Auth User", "auth@example.com", "password123")

		// Test successful authentication
		user, err := userModel.Authenticate("auth@example.com", "password123")
		if err != nil {
			t.Fatalf("Failed to authenticate user: %v", err)
		}

		if user.Email != "auth@example.com" {
			t.Errorf("Expected email 'auth@example.com', got '%s'", user.Email)
		}

		// Test failed authentication
		_, err = userModel.Authenticate("auth@example.com", "wrongpassword")
		if err == nil {
			t.Error("Expected authentication to fail with wrong password")
		}
	})

	// Test email existence check
	t.Run("EmailExists", func(t *testing.T) {
		// Create a user first
		userModel.Create("Existing User", "existing@example.com", "password123")

		// Test existing email
		exists, err := userModel.EmailExists("existing@example.com")
		if err != nil {
			t.Fatalf("Failed to check email existence: %v", err)
		}

		if !exists {
			t.Error("Expected email to exist")
		}

		// Test non-existing email
		exists, err = userModel.EmailExists("nonexisting@example.com")
		if err != nil {
			t.Fatalf("Failed to check email existence: %v", err)
		}

		if exists {
			t.Error("Expected email to not exist")
		}
	})

	// Test get user by ID
	t.Run("GetUserByID", func(t *testing.T) {
		// Create a user first
		createdUser, _ := userModel.Create("Get User", "getuser@example.com", "password123")

		// Get user by ID
		user, err := userModel.GetByID(createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}

		if user.ID != createdUser.ID {
			t.Errorf("Expected ID %d, got %d", createdUser.ID, user.ID)
		}

		if user.Name != "Get User" {
			t.Errorf("Expected name 'Get User', got '%s'", user.Name)
		}
	})
}

// TestBlogModel tests blog model functionality
func TestBlogModel(t *testing.T) {
	setupTestDB()
	defer cleanupTestDB()

	userModel := models.NewUserModel(testDB)
	blogModel := models.NewBlogModel(testDB)

	// Create a test user first
	testUser, err := userModel.Create("Blog Author", "author@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test blog creation
	t.Run("CreateBlog", func(t *testing.T) {
		blog, err := blogModel.Create("Test Blog Title", "This is test blog content.", testUser.ID)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		if blog.Title != "Test Blog Title" {
			t.Errorf("Expected title 'Test Blog Title', got '%s'", blog.Title)
		}

		if blog.Content != "This is test blog content." {
			t.Errorf("Expected content 'This is test blog content.', got '%s'", blog.Content)
		}

		if blog.UserID != testUser.ID {
			t.Errorf("Expected user ID %d, got %d", testUser.ID, blog.UserID)
		}

		if blog.ID == 0 {
			t.Error("Expected blog ID to be set")
		}
	})

	// Test get blog by ID
	t.Run("GetBlogByID", func(t *testing.T) {
		// Create a blog first
		createdBlog, _ := blogModel.Create("Get Blog Test", "Get blog content.", testUser.ID)

		// Get blog by ID
		blog, err := blogModel.GetByID(createdBlog.ID)
		if err != nil {
			t.Fatalf("Failed to get blog by ID: %v", err)
		}

		if blog.ID != createdBlog.ID {
			t.Errorf("Expected ID %d, got %d", createdBlog.ID, blog.ID)
		}

		if blog.Title != "Get Blog Test" {
			t.Errorf("Expected title 'Get Blog Test', got '%s'", blog.Title)
		}

		if blog.UserName != "Blog Author" {
			t.Errorf("Expected user name 'Blog Author', got '%s'", blog.UserName)
		}
	})

	// Test get blogs by user ID
	t.Run("GetBlogsByUserID", func(t *testing.T) {
		// Create multiple blogs for the user
		blogModel.Create("User Blog 1", "Content 1", testUser.ID)
		blogModel.Create("User Blog 2", "Content 2", testUser.ID)

		// Get blogs by user ID
		blogs, err := blogModel.GetByUserID(testUser.ID)
		if err != nil {
			t.Fatalf("Failed to get blogs by user ID: %v", err)
		}

		if len(blogs) < 2 {
			t.Errorf("Expected at least 2 blogs, got %d", len(blogs))
		}

		// Check if all blogs belong to the test user
		for _, blog := range blogs {
			if blog.UserID != testUser.ID {
				t.Errorf("Expected user ID %d, got %d", testUser.ID, blog.UserID)
			}
		}
	})

	// Test blog update
	t.Run("UpdateBlog", func(t *testing.T) {
		// Create a blog first
		createdBlog, _ := blogModel.Create("Original Title", "Original content.", testUser.ID)

		// Update the blog
		updatedBlog, err := blogModel.Update(createdBlog.ID, "Updated Title", "Updated content.")
		if err != nil {
			t.Fatalf("Failed to update blog: %v", err)
		}

		if updatedBlog.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", updatedBlog.Title)
		}

		if updatedBlog.Content != "Updated content." {
			t.Errorf("Expected content 'Updated content.', got '%s'", updatedBlog.Content)
		}
	})

	// Test blog deletion
	t.Run("DeleteBlog", func(t *testing.T) {
		// Create a blog first
		createdBlog, _ := blogModel.Create("Delete Test", "Delete content.", testUser.ID)

		// Delete the blog
		err := blogModel.Delete(createdBlog.ID)
		if err != nil {
			t.Fatalf("Failed to delete blog: %v", err)
		}

		// Try to get the deleted blog (should fail)
		_, err = blogModel.GetByID(createdBlog.ID)
		if err == nil {
			t.Error("Expected error when getting deleted blog")
		}
	})

	// Test user edit permission
	t.Run("CanUserEdit", func(t *testing.T) {
		// Create another user
		anotherUser, _ := userModel.Create("Another User", "another@example.com", "password123")

		// Create a blog by test user
		blog, _ := blogModel.Create("Permission Test", "Permission content.", testUser.ID)

		// Test that owner can edit
		canEdit, err := blogModel.CanUserEdit(blog.ID, testUser.ID)
		if err != nil {
			t.Fatalf("Failed to check edit permission: %v", err)
		}

		if !canEdit {
			t.Error("Expected owner to be able to edit")
		}

		// Test that non-owner cannot edit
		canEdit, err = blogModel.CanUserEdit(blog.ID, anotherUser.ID)
		if err != nil {
			t.Fatalf("Failed to check edit permission: %v", err)
		}

		if canEdit {
			t.Error("Expected non-owner to not be able to edit")
		}
	})
}

// TestBlogCount tests blog counting functionality
func TestBlogCount(t *testing.T) {
	setupTestDB()
	defer cleanupTestDB()

	userModel := models.NewUserModel(testDB)
	blogModel := models.NewBlogModel(testDB)

	// Create a test user
	testUser, _ := userModel.Create("Count User", "count@example.com", "password123")

	// Initially, there should be 0 blogs
	count, err := blogModel.Count()
	if err != nil {
		t.Fatalf("Failed to count blogs: %v", err)
	}

	initialCount := count

	// Create some blogs
	blogModel.Create("Count Blog 1", "Content 1", testUser.ID)
	blogModel.Create("Count Blog 2", "Content 2", testUser.ID)
	blogModel.Create("Count Blog 3", "Content 3", testUser.ID)

	// Count should increase by 3
	count, err = blogModel.Count()
	if err != nil {
		t.Fatalf("Failed to count blogs: %v", err)
	}

	expectedCount := initialCount + 3
	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}
}
