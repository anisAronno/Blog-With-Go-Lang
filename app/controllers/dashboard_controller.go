// app/controllers/dashboard_controller.go - Handles dashboard related pages
package controllers

import (
	"go-web-app/app/middleware"
	"go-web-app/app/models"
	"go-web-app/config"
	"net/http"
	"strconv"
	"strings"

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

	// Get user's recent blogs
	userBlogs, err := c.BlogModel.GetByUserID(user.ID)
	if err != nil {
		userBlogs = []*models.Blog{} // Default to empty slice on error
	}

	// Get user's blog statistics
	totalUserBlogs := len(userBlogs)
	publishedUserBlogs, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "published")
	draftUserBlogs, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "draft")

	// Stats structure for template
	stats := map[string]interface{}{
		"TotalBlogs":     totalUserBlogs,
		"PublishedBlogs": publishedUserBlogs,
		"DraftBlogs":     draftUserBlogs,
	}

	// For admin users, add global statistics
	if user.IsAdmin() {
		totalUsers, _ := c.UserModel.Count()
		totalBlogs, _ := c.BlogModel.Count()
		totalPublished, _ := c.BlogModel.CountByStatus("published")
		totalDrafts, _ := c.BlogModel.CountByStatus("draft")
		totalAdmins, _ := c.UserModel.CountByRole("admin")
		totalAuthors, _ := c.UserModel.CountByRole("author")
		totalRegularUsers, _ := c.UserModel.CountByRole("user")

		stats["TotalUsers"] = totalUsers
		stats["GlobalTotalBlogs"] = totalBlogs
		stats["GlobalPublishedBlogs"] = totalPublished
		stats["GlobalDraftBlogs"] = totalDrafts
		stats["TotalAdmins"] = totalAdmins
		stats["TotalAuthors"] = totalAuthors
		stats["TotalRegularUsers"] = totalRegularUsers
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":       "Dashboard",
		"User":        user,
		"Stats":       stats,
		"RecentBlogs": userBlogs[:min(len(userBlogs), 5)], // Show last 5 blogs
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

	// Get user's blog count and stats
	userBlogs, err := c.BlogModel.GetByUserID(user.ID)
	if err != nil {
		userBlogs = []*models.Blog{}
	}

	publishedCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "published")
	draftCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "draft")

	userStats := map[string]interface{}{
		"TotalBlogs":     len(userBlogs),
		"PublishedBlogs": publishedCount,
		"DraftBlogs":     draftCount,
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":     "My Profile",
		"User":      user,
		"BlogCount": len(userBlogs),
		"UserStats": userStats,
	}

	renderTemplate(w, "dashboard/profile", data)
}

// ChangePassword handles password change for the current user
func (c *DashboardController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Helper function to show profile with error
	showProfileWithError := func(errorMsg string) {
		userBlogs, _ := c.BlogModel.GetByUserID(user.ID)
		publishedCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "published")
		draftCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "draft")

		userStats := map[string]interface{}{
			"TotalBlogs":     len(userBlogs),
			"PublishedBlogs": publishedCount,
			"DraftBlogs":     draftCount,
		}

		data := map[string]interface{}{
			"Title":     "My Profile",
			"User":      user,
			"BlogCount": len(userBlogs),
			"UserStats": userStats,
			"Error":     errorMsg,
		}
		renderTemplate(w, "dashboard/profile", data)
	}

	// Get form data
	currentPassword := strings.TrimSpace(r.FormValue("current_password"))
	newPassword := strings.TrimSpace(r.FormValue("new_password"))
	confirmPassword := strings.TrimSpace(r.FormValue("confirm_password"))

	// Validate form data
	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		showProfileWithError("All password fields are required")
		return
	}

	if newPassword != confirmPassword {
		showProfileWithError("New password and confirmation do not match")
		return
	}

	if len(newPassword) < 6 {
		showProfileWithError("New password must be at least 6 characters long")
		return
	}

	// Change password
	err = c.UserModel.UpdateProfile(user.ID, user.Name, user.Email, currentPassword, &newPassword)
	if err != nil {
		if strings.Contains(err.Error(), "current password is incorrect") {
			showProfileWithError("Current password is incorrect")
			return
		}
		showProfileWithError("Failed to change password: " + err.Error())
		return
	}

	// Show success message
	userBlogs, _ := c.BlogModel.GetByUserID(user.ID)
	publishedCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "published")
	draftCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "draft")

	userStats := map[string]interface{}{
		"TotalBlogs":     len(userBlogs),
		"PublishedBlogs": publishedCount,
		"DraftBlogs":     draftCount,
	}

	data := map[string]interface{}{
		"Title":     "My Profile",
		"User":      user,
		"BlogCount": len(userBlogs),
		"UserStats": userStats,
		"Success":   "Password changed successfully",
	}
	renderTemplate(w, "dashboard/profile", data)
}

// UpdateProfile updates the current user's profile
// UpdateProfile updates the current user's profile (name and email only)
func (c *DashboardController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Helper function to show profile form with error
	showProfileWithError := func(errorMsg string) {
		userBlogs, _ := c.BlogModel.GetByUserID(user.ID)
		publishedCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "published")
		draftCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "draft")

		userStats := map[string]interface{}{
			"TotalBlogs":     len(userBlogs),
			"PublishedBlogs": publishedCount,
			"DraftBlogs":     draftCount,
		}

		data := map[string]interface{}{
			"Title":     "My Profile",
			"User":      user,
			"BlogCount": len(userBlogs),
			"UserStats": userStats,
			"Error":     errorMsg,
		}
		renderTemplate(w, "dashboard/profile", data)
	}

	// Get form data
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))

	// Validate required fields
	if name == "" || email == "" {
		showProfileWithError("Name and email are required")
		return
	}

	// Validate email format
	if !strings.Contains(email, "@") {
		showProfileWithError("Invalid email format")
		return
	}

	// Update profile (name and email only, without password change)
	err = c.UserModel.UpdateProfile(user.ID, name, email, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			showProfileWithError("Email address is already in use by another user")
			return
		}
		showProfileWithError("Failed to update profile: " + err.Error())
		return
	}

	// Get updated user data
	updatedUser, err := c.UserModel.GetByID(user.ID)
	if err != nil {
		updatedUser = user // fallback to current user if error
	}

	// Get user stats
	userBlogs, _ := c.BlogModel.GetByUserID(updatedUser.ID)
	publishedCount, _ := c.BlogModel.CountUserBlogsByStatus(updatedUser.ID, "published")
	draftCount, _ := c.BlogModel.CountUserBlogsByStatus(updatedUser.ID, "draft")

	userStats := map[string]interface{}{
		"TotalBlogs":     len(userBlogs),
		"PublishedBlogs": publishedCount,
		"DraftBlogs":     draftCount,
	}

	// Show success message
	data := map[string]interface{}{
		"Title":     "My Profile",
		"User":      updatedUser,
		"BlogCount": len(userBlogs),
		"UserStats": userStats,
		"Success":   "Profile updated successfully",
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

	// Calculate user role statistics
	adminCount, _ := c.UserModel.CountByRole("admin")
	authorCount, _ := c.UserModel.CountByRole("author")
	userCount, _ := c.UserModel.CountByRole("user")

	// Prepare user stats
	userStats := map[string]interface{}{
		"AdminCount":  adminCount,
		"AuthorCount": authorCount,
		"UserCount":   userCount,
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":       "All Users",
		"Users":       users,
		"User":        currentUser,
		"CurrentUser": currentUser,
		"IsAdmin":     currentUser.IsAdmin(),
		"UserStats":   userStats,
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

	// Prevent deletion of super admin (ID 1)
	if userID == 1 {
		http.Error(w, "Super admin account cannot be deleted", http.StatusForbidden)
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

// EditUser shows the user edit form (admin only)
func (c *DashboardController) EditUser(w http.ResponseWriter, r *http.Request) {
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
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prevent editing super admin (ID 1)
	if userID == 1 {
		http.Error(w, "Super admin account cannot be edited", http.StatusForbidden)
		return
	}

	// Get user to edit
	editUser, err := c.UserModel.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":    "Edit User",
		"User":     currentUser,
		"EditUser": editUser,
	}

	renderTemplate(w, "dashboard/users/edit", data)
}

// UpdateUser updates user information (admin only)
func (c *DashboardController) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prevent editing super admin (ID 1)
	if userID == 1 {
		http.Error(w, "Super admin account cannot be edited", http.StatusForbidden)
		return
	}

	// Get user to edit for error display
	editUser, err := c.UserModel.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Helper function to show edit form with error
	showEditWithError := func(errorMsg string) {
		data := map[string]interface{}{
			"Title":    "Edit User",
			"User":     currentUser,
			"EditUser": editUser,
			"Error":    errorMsg,
		}
		renderTemplate(w, "dashboard/users/edit", data)
	}

	// Get form data
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	role := strings.TrimSpace(r.FormValue("role"))
	password := strings.TrimSpace(r.FormValue("password"))

	// Validate required fields
	if name == "" || email == "" || role == "" {
		showEditWithError("Name, email, and role are required")
		return
	}

	// Validate role
	if role != "user" && role != "author" && role != "admin" {
		showEditWithError("Invalid role selected")
		return
	}

	// Validate email format
	if !strings.Contains(email, "@") {
		showEditWithError("Invalid email format")
		return
	}

	// Update user
	var passwordPtr *string
	if password != "" {
		passwordPtr = &password
	}

	err = c.UserModel.Update(userID, name, email, role, passwordPtr)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			showEditWithError("Email address is already in use by another user")
			return
		}
		showEditWithError("Failed to update user: " + err.Error())
		return
	}

	// Redirect to users list after successful update
	http.Redirect(w, r, "/dashboard/users", http.StatusSeeOther)
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
