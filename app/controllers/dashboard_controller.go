// app/controllers/dashboard_controller.go - Handles dashboard related pages
package controllers

import (
	"go-web-app/app/middleware"
	"go-web-app/app/models"
	"go-web-app/config"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DashboardController handles dashboard pages
type DashboardController struct {
	UserModel *models.UserModel
	BlogModel *models.BlogModel
}

// NewDashboardController creates a new DashboardController
func NewDashboardController() *DashboardController {
	return &DashboardController{
		UserModel: models.NewUserModel(config.Database),
		BlogModel: models.NewBlogModel(config.Database),
	}
}

// Index displays the main dashboard
func (c *DashboardController) Index(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Get user's blog count
	userBlogs, err := c.BlogModel.GetByUserID(user.ID)
	if err != nil {
		userBlogs = []*models.Blog{} // Default to empty slice on error
	}

	// Get total blog count
	totalBlogs, err := c.BlogModel.Count()
	if err != nil {
		totalBlogs = 0
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":         "Dashboard",
		"User":          user,
		"UserBlogCount": len(userBlogs),
		"TotalBlogs":    totalBlogs,
		"RecentBlogs":   userBlogs[:min(len(userBlogs), 5)], // Show last 5 blogs
	}

	renderTemplate(w, "dashboard/index", data)
}

// Profile displays the user profile page
func (c *DashboardController) Profile(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Get user's blog count
	userBlogs, err := c.BlogModel.GetByUserID(user.ID)
	if err != nil {
		userBlogs = []*models.Blog{}
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":     "My Profile",
		"User":      user,
		"BlogCount": len(userBlogs),
	}

	renderTemplate(w, "dashboard/profile", data)
}

// Users displays all users (admin only)
func (c *DashboardController) Users(w http.ResponseWriter, r *http.Request) {
	// Get current user
	currentUser, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user is admin
	if !currentUser.IsAdmin() {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Get all users
	users, err := c.UserModel.GetAll()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":       "All Users",
		"Users":       users,
		"User":        currentUser,
		"CurrentUser": currentUser,
		"IsAdmin":     currentUser.IsAdmin(),
	}

	renderTemplate(w, "dashboard/users", data)
}

// DeleteUser deletes a user (admin only)
func (c *DashboardController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard/users", http.StatusSeeOther)
		return
	}

	// Get current user
	currentUser, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if current user is admin
	if !currentUser.IsAdmin() {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Get user ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prevent admin from deleting themselves
	if userID == currentUser.ID {
		http.Error(w, "You cannot delete your own account", http.StatusBadRequest)
		return
	}

	// Delete user
	err = c.UserModel.Delete(userID)
	if err != nil {
		// Check if this is the main admin protection error
		if err.Error() == "cannot delete the main administrator account" {
			http.Error(w, "Cannot delete the main administrator account", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Redirect to users list
	http.Redirect(w, r, "/dashboard/users", http.StatusSeeOther)
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
