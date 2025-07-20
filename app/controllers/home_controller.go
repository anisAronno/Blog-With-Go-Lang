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
	// Get page parameter from URL (default to 1)
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Calculate offset for pagination
	limit := 12
	offset := (page - 1) * limit

	// Get blogs for current page
	blogs, err := c.BlogModel.GetAll(limit, offset)
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
		return
	}

	// Get total blog count for pagination
	totalBlogs, err := c.BlogModel.Count()
	if err != nil {
		totalBlogs = 0 // Default to 0 if count fails
	}

	// Calculate pagination info
	totalPages := (totalBlogs + limit - 1) / limit // Ceiling division
	hasNext := page < totalPages
	hasPrev := page > 1

	// Check if user is logged in (using session for homepage)
	user, _ := middleware.GetCurrentUserFromSession(r) // Don't show error for homepage

	// Prepare data for template
	data := map[string]interface{}{
		"Title":      "Welcome to Go Blog",
		"Blogs":      blogs,
		"User":       user, // Include user for navigation
		"Page":       page,
		"TotalPages": totalPages,
		"HasNext":    hasNext,
		"HasPrev":    hasPrev,
		"NextPage":   page + 1,
		"PrevPage":   page - 1,
		"BaseURL":    "/", // For pagination component
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
	user, _ := middleware.GetCurrentUserFromSession(r)

	// Prepare data for template
	data := map[string]interface{}{
		"Title": blog.Title,
		"Blog":  blog,
		"User":  user, // Add user data to template
	}

	renderTemplate(w, "blog/show", data)
}
