// database/seeders/blog_seeder.go - Seed blog data for testing
package main

import (
	"go-web-app/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	appConfig := config.LoadConfig()
	
	// Connect to database
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	log.Println("🌱 Seeding blog data...")

	// Get user IDs to assign blogs to
	var userIDs []int
	userQuery := `SELECT id FROM users ORDER BY id`
	rows, err := db.Query(userQuery)
	if err != nil {
		log.Fatal("Failed to get users: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			log.Fatal("Failed to scan user ID: ", err)
		}
		userIDs = append(userIDs, userID)
	}

	if len(userIDs) == 0 {
		log.Fatal("No users found. Please create users first.")
	}

	// Sample blog data
	blogs := []struct {
		Title   string
		Content string
		UserID  int
	}{
		{
			"Getting Started with Go Web Development",
			"Go is an excellent choice for web development. Its simplicity, performance, and built-in concurrency make it perfect for building scalable web applications. In this post, we'll explore the fundamentals of web development with Go, including routing, middleware, and template rendering.",
			userIDs[0],
		},
		{
			"Building RESTful APIs with Go",
			"REST APIs are the backbone of modern web applications. Go provides excellent tools for building robust, scalable APIs. We'll cover HTTP routing, JSON handling, middleware implementation, and best practices for API design and security.",
			userIDs[0],
		},
		{
			"Database Integration in Go",
			"Working with databases is crucial for most web applications. This guide covers connecting to MySQL, performing CRUD operations, handling transactions, and implementing proper error handling. We'll also discuss connection pooling and performance optimization.",
			userIDs[len(userIDs)-1],
		},
		{
			"Go Templates and Frontend Integration",
			"Go's template package provides powerful tools for generating dynamic HTML. Learn how to create reusable templates, handle template inheritance, pass data efficiently, and integrate with modern frontend frameworks.",
			userIDs[0],
		},
		{
			"Authentication and Authorization in Go",
			"Security is paramount in web applications. This comprehensive guide covers implementing user authentication, session management, password hashing with bcrypt, role-based access control, and protecting against common security vulnerabilities.",
			userIDs[len(userIDs)-1],
		},
		{
			"Deployment and Production Best Practices",
			"Taking your Go web application to production requires careful planning. We'll discuss Docker containerization, environment configuration, monitoring, logging, performance tuning, and scaling strategies for high-traffic applications.",
			userIDs[0],
		},
		{
			"Testing Go Web Applications",
			"Quality assurance through testing is essential for reliable applications. Learn about unit testing, integration testing, mocking dependencies, testing HTTP handlers, and implementing continuous integration pipelines.",
			userIDs[len(userIDs)-1],
		},
		{
			"Microservices Architecture with Go",
			"Go's performance and simplicity make it ideal for microservices. Explore service decomposition, inter-service communication, distributed tracing, service discovery, and handling the complexities of distributed systems.",
			userIDs[0],
		},
	}

	// Insert sample blogs
	insertQuery := `INSERT INTO blogs (title, content, user_id, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())`
	
	for _, blog := range blogs {
		_, err := db.Exec(insertQuery, blog.Title, blog.Content, blog.UserID)
		if err != nil {
			log.Printf("Warning: Failed to insert blog '%s': %v", blog.Title, err)
		} else {
			log.Printf("✅ Created blog: %s", blog.Title)
		}
	}

	log.Println("🎉 Blog seeding completed successfully!")
}
