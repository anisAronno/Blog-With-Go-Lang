package seeders

import (
	"database/sql"
	"fmt"
	"go-web-app/app/models"
	"log"
	"math/rand"
	"time"
)

// BlogSeeder handles blog data seeding
type BlogSeeder struct {
	DB        *sql.DB
	BlogModel *models.BlogModel
	UserModel *models.UserModel
}

// NewBlogSeeder creates a new blog seeder
func NewBlogSeeder(db *sql.DB) *BlogSeeder {
	return &BlogSeeder{
		DB:        db,
		BlogModel: models.NewBlogModel(db),
		UserModel: models.NewUserModel(db),
	}
}

// Seed creates sample blog posts
func (s *BlogSeeder) Seed() error {
	log.Println("🌱 Seeding blogs...")

	// Get all users to assign blogs to
	users, err := s.UserModel.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	if len(users) == 0 {
		log.Println("⚠️  No users found, skipping blog seeding")
		return nil
	}

	blogTitles := []string{
		"Getting Started with Go Programming",
		"Understanding Database Design Patterns",
		"Modern Web Development with Go",
		"Building RESTful APIs",
		"Docker and Container Orchestration",
		"Introduction to Machine Learning",
		"Cloud Computing Best Practices",
		"Cybersecurity Fundamentals",
		"Mobile App Development Trends",
		"DevOps Culture and Practices",
	}

	blogContents := []string{
		"Go is a powerful programming language developed by Google. It offers excellent performance and simplicity that makes it perfect for modern web applications.",
		"Database design patterns are crucial for building scalable applications. Understanding normalization, indexing, and relationships is key to success.",
		"Modern web development requires understanding of both frontend and backend technologies. Go provides excellent tools for building robust web applications.",
		"RESTful APIs are the backbone of modern web services. They provide a standardized way for applications to communicate with each other.",
		"Docker has revolutionized how we deploy and manage applications. Container orchestration with Kubernetes takes this to the next level.",
		"Machine learning is transforming industries. Understanding the basics of ML algorithms and data processing is becoming essential for developers.",
		"Cloud computing offers scalability and flexibility. Learning best practices for cloud deployment ensures reliable and cost-effective solutions.",
		"Cybersecurity is more important than ever. Understanding common vulnerabilities and security practices protects both users and businesses.",
		"Mobile app development continues to evolve. Cross-platform frameworks and native development each have their place in modern app development.",
		"DevOps culture emphasizes collaboration between development and operations teams. Automation and continuous integration are key practices.",
	}

	// Seed 15 blog posts
	for i := 0; i < 15; i++ {
		// Randomly select a user and blog content
		userIndex := rand.Intn(len(users))
		titleIndex := rand.Intn(len(blogTitles))
		contentIndex := rand.Intn(len(blogContents))

		user := users[userIndex]
		title := fmt.Sprintf("%s - Part %d", blogTitles[titleIndex], (i%3)+1)
		content := fmt.Sprintf("%s\n\nThis is a detailed exploration of the topic with practical examples and real-world applications. The content provides valuable insights for developers and technology enthusiasts.", blogContents[contentIndex])

		// Check if similar blog already exists
		query := `SELECT COUNT(*) FROM blogs WHERE title = ? AND user_id = ?`
		var count int
		err := s.DB.QueryRow(query, title, user.ID).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check existing blog: %v", err)
		}

		if count > 0 {
			log.Printf("⏭️  Blog '%s' already exists, skipping", title)
			continue
		}

		// Create blog
		_, err = s.BlogModel.Create(title, content, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create blog: %v", err)
		}

		log.Printf("✅ Created blog: '%s' by %s", title, user.Name)
	}

	log.Println("✅ Blog seeding completed")
	return nil
}

// Clear removes all blog posts
func (s *BlogSeeder) Clear() error {
	log.Println("🧹 Clearing blog data...")

	query := `DELETE FROM blogs`
	_, err := s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear blogs: %v", err)
	}

	log.Println("✅ Blog data cleared")
	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
