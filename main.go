// main.go - Entry point of the Go Web Application
// This application follows Laravel-like structure for easy learning and maintenance
package main

import (
	"fmt"
	"go-web-app/app/middleware"
	"go-web-app/config"
	"go-web-app/routes"
	"log"
	"net/http"
)

// main function bootstraps the application
func main() {
	// 1. Load application configuration from .env file
	appConfig := config.LoadConfig()
	fmt.Printf("🚀 Starting Go Web App in %s mode\n", appConfig.AppEnv)

	// 2. Connect to MySQL database
	db, err := config.ConnectDatabase(appConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// 3. Initialize sessions for user authentication
	middleware.InitSessions()
	fmt.Println("✅ Sessions initialized")

	// 4. Setup application routes (similar to Laravel's web.php)
	router := routes.SetupRoutes()

	// 5. Apply global middleware
	handler := middleware.LoggingMiddleware(
		middleware.CORSMiddleware(router),
	)

	// 6. Start the HTTP server
	fmt.Printf("🌐 Server started at http://localhost:%s\n", appConfig.AppPort)
	fmt.Println("📝 Available routes:")
	fmt.Println("   - GET  /                 (Homepage - Blog listing)")
	fmt.Println("   - GET  /login            (Login page)")
	fmt.Println("   - GET  /register         (Register page)")
	fmt.Println("   - GET  /dashboard        (Dashboard - requires auth)")
	fmt.Println("   - GET  /dashboard/blogs  (Blog management)")
	fmt.Println("   - GET  /dashboard/admin/blogs (Admin blog management)")
	fmt.Println("   - GET  /dashboard/users  (User listing)")
	fmt.Println("   - GET  /dashboard/profile (User profile)")

	err = http.ListenAndServe(":"+appConfig.AppPort, handler)
	if err != nil {
		log.Fatal("❌ Server error: ", err)
	}
}
