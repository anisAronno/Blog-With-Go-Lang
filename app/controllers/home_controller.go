// app/controllers/home_controller.go - Handles public homepage and blog display
package controllers

import (
	"go-web-app/app/middleware"
	"go-web-app/app/models"
	"go-web-app/config"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HomeController handles public pages
type HomeController struct {
	BlogModel *models.BlogModel
	UserModel *models.UserModel
}

// NewHomeController creates a new HomeController
func NewHomeController() *HomeController {
	return &HomeController{
		BlogModel: models.NewBlogModel(config.Database),
		UserModel: models.NewUserModel(config.Database),
	}
}

// Index displays the homepage with blog listing
func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	// Get all blogs for homepage
	blogs, err := c.BlogModel.GetAll(10, 0) // Get latest 10 blogs
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
		return
	}

	// Check if user is logged in (using session for homepage)
	user, _ := middleware.GetCurrentUserFromSession(r) // Don't show error for homepage

	// Prepare data for template
	data := map[string]interface{}{
		"Title": "Welcome to Go Blog",
		"Blogs": blogs,
		"User":  user, // Include user for navigation
	}

	renderTemplate(w, "home", data)
}

// ShowBlog displays a single blog post
func (c *HomeController) ShowBlog(w http.ResponseWriter, r *http.Request) {
	// Get blog ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	// Get blog by ID
	blog, err := c.BlogModel.GetByID(id)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Get current user (if logged in)
	user, _ := middleware.GetCurrentUser(r)

	// Prepare data for template
	data := map[string]interface{}{
		"Title": blog.Title,
		"Blog":  blog,
		"User":  user, // Add user data to template
	}

	renderTemplate(w, "blog/show", data)
}
