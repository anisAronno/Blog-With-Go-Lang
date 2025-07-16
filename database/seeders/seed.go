// database/seeders/seed.go - Database seeder for initial data
package main

import (
	"go-web-app/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration and connect to database
	appConfig := config.LoadConfig()
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	log.Println("🌱 Seeding database with initial data...")

	// Hash password for test users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password: ", err)
	}

	// Insert test users
	users := []map[string]string{
		{"name": "Admin User", "email": "admin@example.com"},
		{"name": "John Doe", "email": "john@example.com"},
		{"name": "Jane Smith", "email": "jane@example.com"},
	}

	for _, user := range users {
		// Check if user already exists
		var count int
		checkQuery := "SELECT COUNT(*) FROM users WHERE email = ?"
		err = db.QueryRow(checkQuery, user["email"]).Scan(&count)
		if err != nil {
			log.Printf("Error checking user %s: %v", user["email"], err)
			continue
		}

		if count == 0 {
			insertQuery := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
			_, err = db.Exec(insertQuery, user["name"], user["email"], string(hashedPassword))
			if err != nil {
				log.Printf("Failed to create user %s: %v", user["email"], err)
			} else {
				log.Printf("✅ Created user: %s (%s)", user["name"], user["email"])
			}
		} else {
			log.Printf("👤 User already exists: %s", user["email"])
		}
	}

	// Insert sample blog posts
	blogPosts := []map[string]interface{}{
		{
			"title":   "Welcome to Go Web Application",
			"content": "This is a comprehensive Go web application built with Laravel-like structure. It demonstrates CRUD operations, authentication, and modern web development practices using Go.\n\nFeatures include:\n- User authentication (login/register/logout)\n- Blog post management (CRUD operations)\n- Dashboard with beautiful Tailwind CSS design\n- Session management\n- MySQL database integration\n- Enterprise-grade project structure",
			"email":   "admin@example.com",
		},
		{
			"title":   "Getting Started with Go Web Development",
			"content": "Go (Golang) is an excellent choice for web development due to its simplicity, performance, and built-in concurrency support. This application showcases how to build a modern web application using Go with patterns familiar to Laravel developers.\n\nKey concepts covered:\n- MVC architecture\n- Database migrations and seeders\n- Middleware for authentication and logging\n- Template rendering with HTML templates\n- Form handling and validation",
			"email":   "john@example.com",
		},
		{
			"title":   "Building Scalable Web Applications",
			"content": "When building web applications, it's important to follow best practices that allow for easy maintenance and scaling. This project demonstrates:\n\n- Clean separation of concerns\n- Modular architecture\n- Environment-based configuration\n- Comprehensive error handling\n- Security best practices\n\nThe folder structure follows Laravel conventions, making it easy for PHP developers to understand and contribute to Go projects.",
			"email":   "jane@example.com",
		},
		{
			"title":   "Database Design and Relationships",
			"content": "This application uses a simple but effective database design with two main tables: users and blogs. The relationship between users and blogs is a one-to-many relationship, where each user can have multiple blog posts.\n\nThe database schema includes:\n- Proper indexing for performance\n- Foreign key constraints for data integrity\n- Timestamps for audit trails\n- UTF8MB4 charset for full Unicode support\n\nThis design can be easily extended to include additional features like categories, tags, comments, and more.",
			"email":   "admin@example.com",
		},
		{
			"title":   "Frontend Design with Tailwind CSS",
			"content": "The frontend of this application is built using Tailwind CSS, a utility-first CSS framework that enables rapid UI development. The design focuses on:\n\n- Clean and modern interface\n- Responsive design for all devices\n- Consistent color scheme and typography\n- Intuitive navigation and user experience\n- Accessibility considerations\n\nThe dashboard provides an easy-to-use interface for managing blog posts and viewing user information. The public homepage showcases blog posts in an attractive, readable format.",
			"email":   "jane@example.com",
		},
	}

	for _, post := range blogPosts {
		// Get user ID by email
		var userID int
		getUserQuery := "SELECT id FROM users WHERE email = ?"
		err = db.QueryRow(getUserQuery, post["email"]).Scan(&userID)
		if err != nil {
			log.Printf("Failed to find user with email %s: %v", post["email"], err)
			continue
		}

		// Check if blog post already exists
		var count int
		checkQuery := "SELECT COUNT(*) FROM blogs WHERE title = ? AND user_id = ?"
		err = db.QueryRow(checkQuery, post["title"], userID).Scan(&count)
		if err != nil {
			log.Printf("Error checking blog post %s: %v", post["title"], err)
			continue
		}

		if count == 0 {
			insertQuery := "INSERT INTO blogs (title, content, user_id, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
			_, err = db.Exec(insertQuery, post["title"], post["content"], userID)
			if err != nil {
				log.Printf("Failed to create blog post %s: %v", post["title"], err)
			} else {
				log.Printf("✅ Created blog post: %s", post["title"])
			}
		} else {
			log.Printf("📝 Blog post already exists: %s", post["title"])
		}
	}

	log.Println("🎉 Database seeding completed successfully!")
	log.Println("")
	log.Println("👤 Test Users Created:")
	log.Println("   Email: admin@example.com | Password: password")
	log.Println("   Email: john@example.com  | Password: password")
	log.Println("   Email: jane@example.com  | Password: password")
}
