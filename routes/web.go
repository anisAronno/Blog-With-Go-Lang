// routes/web.go - Application routes (similar to Laravel's web.php)
package routes

import (
	"go-web-app/app/controllers"
	"go-web-app/app/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes() *mux.Router {
	// Initialize router
	r := mux.NewRouter()

	// Initialize controllers
	authController := controllers.NewAuthController()
	homeController := controllers.NewHomeController()
	dashboardController := controllers.NewDashboardController()
	blogController := controllers.NewBlogController()

	// Static files serving
	r.PathPrefix("/public/").Handler(controllers.StaticFileHandler())

	// Public routes (accessible to everyone)
	r.HandleFunc("/", homeController.Index).Methods("GET")
	r.HandleFunc("/blog/{id}", homeController.ShowBlog).Methods("GET")

	// Guest routes (only for non-authenticated users)
	r.HandleFunc("/login", middleware.GuestMiddleware(authController.ShowLogin)).Methods("GET")
	r.HandleFunc("/login", middleware.GuestMiddleware(authController.Login)).Methods("POST")
	r.HandleFunc("/register", middleware.GuestMiddleware(authController.ShowRegister)).Methods("GET")
	r.HandleFunc("/register", middleware.GuestMiddleware(authController.Register)).Methods("POST")

	// Authentication route
	r.HandleFunc("/logout", authController.Logout).Methods("POST")

	// Protected routes (require authentication)
	// Dashboard routes
	dashboard := r.PathPrefix("/dashboard").Subrouter()
	dashboard.HandleFunc("", middleware.AuthMiddleware(dashboardController.Index)).Methods("GET")
	dashboard.HandleFunc("/", middleware.AuthMiddleware(dashboardController.Index)).Methods("GET")
	dashboard.HandleFunc("/profile", middleware.AuthMiddleware(dashboardController.Profile)).Methods("GET")
	dashboard.HandleFunc("/users", middleware.AuthMiddleware(dashboardController.Users)).Methods("GET")
	dashboard.HandleFunc("/users/{id}/delete", middleware.AuthMiddleware(dashboardController.DeleteUser)).Methods("POST")

	// Blog management routes
	dashboard.HandleFunc("/blogs", middleware.AuthMiddleware(blogController.Index)).Methods("GET")
	dashboard.HandleFunc("/blogs/create", middleware.AuthMiddleware(blogController.Create)).Methods("GET")
	dashboard.HandleFunc("/blogs", middleware.AuthMiddleware(blogController.Store)).Methods("POST")
	dashboard.HandleFunc("/blogs/{id}/edit", middleware.AuthMiddleware(blogController.Edit)).Methods("GET")
	dashboard.HandleFunc("/blogs/{id}", middleware.AuthMiddleware(blogController.Update)).Methods("POST")
	dashboard.HandleFunc("/blogs/{id}/delete", middleware.AuthMiddleware(blogController.Delete)).Methods("POST")

	// Admin-only blog management routes
	dashboard.HandleFunc("/admin/blogs", middleware.AuthMiddleware(blogController.AdminIndex)).Methods("GET")
	dashboard.HandleFunc("/admin/blogs/{id}/delete", middleware.AuthMiddleware(blogController.AdminDelete)).Methods("POST")

	return r
}
