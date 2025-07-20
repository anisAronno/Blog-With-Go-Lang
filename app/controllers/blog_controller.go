// app/controllers/blog_controller.go - Handles blog CRUD operations
package controllers

import (
	"go-web-app/app/middleware"
	"go-web-app/app/models"
	"go-web-app/config"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// BlogController handles blog CRUD operations
type BlogController struct {
	BlogModel *models.BlogModel
	UserModel *models.UserModel
}

// NewBlogController creates a new BlogController
func NewBlogController() *BlogController {
	return &BlogController{
		BlogModel: models.NewBlogModel(config.Database),
		UserModel: models.NewUserModel(config.Database),
	}
}

// Index displays only the current user's blogs (for personal dashboard)
func (c *BlogController) Index(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Get page parameter from URL (default to 1)
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Calculate offset for pagination
	limit := 10
	offset := (page - 1) * limit

	// Get user's blogs for current page
	blogs, err := c.BlogModel.GetByUserIDPaginated(user.ID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
		return
	}

	// Get total user blog count for pagination
	totalBlogs, err := c.BlogModel.CountUserBlogs(user.ID)
	if err != nil {
		totalBlogs = 0 // Default to 0 if count fails
	}

	// Calculate pagination info
	totalPages := (totalBlogs + limit - 1) / limit // Ceiling division
	hasNext := page < totalPages
	hasPrev := page > 1

	// Debug: Log how many blogs we found
	log.Printf("Blog Controller: Found %d blogs for user %s (page %d)", len(blogs), user.Name, page)

	// Calculate user's blog statistics
	totalUserBlogs := totalBlogs
	publishedCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "published")
	draftCount, _ := c.BlogModel.CountUserBlogsByStatus(user.ID, "draft")

	// Prepare stats
	stats := map[string]interface{}{
		"TotalBlogs":     totalBlogs,
		"PublishedBlogs": publishedCount,
		"DraftBlogs":     draftCount,
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":      "My Blogs",
		"Blogs":      blogs,
		"User":       user,
		"IsAdmin":    user.IsAdmin(),
		"Stats":      stats,
		"Page":       page,
		"TotalPages": totalPages,
		"HasNext":    hasNext,
		"HasPrev":    hasPrev,
		"NextPage":   page + 1,
		"PrevPage":   page - 1,
	}

	renderTemplate(w, "dashboard/blogs/index", data)
}

// AdminIndex displays all blogs for admin users only
func (c *BlogController) AdminIndex(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user is admin
	if !user.IsAdmin() {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Get page parameter from URL (default to 1)
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Calculate offset for pagination
	limit := 10
	offset := (page - 1) * limit

	// Get blogs for current page
	blogs, err := c.BlogModel.GetAllBlogs(limit, offset)
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

	// Debug: Log how many blogs we found
	log.Printf("Admin Blog Controller: Found %d total blogs for admin %s (page %d)", len(blogs), user.Name, page)

	// Calculate global blog statistics for admin
	globalTotalBlogs := totalBlogs
	publishedCount, _ := c.BlogModel.CountByStatus("published")
	draftCount, _ := c.BlogModel.CountByStatus("draft")
	totalAuthors, _ := c.UserModel.CountByRole("author")

	// Prepare stats
	stats := map[string]interface{}{
		"TotalBlogs":     totalBlogs,
		"PublishedBlogs": publishedCount,
		"DraftBlogs":     draftCount,
		"TotalAuthors":   totalAuthors,
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title":      "All Blogs (Admin Management)",
		"Blogs":      blogs,
		"User":       user,
		"IsAdmin":    true,
		"AdminView":  true, // Flag to indicate this is admin view
		"Stats":      stats,
		"Page":       page,
		"TotalPages": totalPages,
		"HasNext":    hasNext,
		"HasPrev":    hasPrev,
		"NextPage":   page + 1,
		"PrevPage":   page - 1,
	}

	renderTemplate(w, "dashboard/blogs/admin", data)
}

// Create shows the create blog form
func (c *BlogController) Create(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title": "Create New Blog",
		"User":  user,
	}

	renderTemplate(w, "dashboard/blogs/create", data)
}

// Store creates a new blog post
func (c *BlogController) Store(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard/blogs/create", http.StatusSeeOther)
		return
	}

	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Get form data
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	excerpt := strings.TrimSpace(r.FormValue("excerpt"))
	status := strings.TrimSpace(r.FormValue("status"))

	// Set defaults
	if excerpt == "" {
		// Auto-generate excerpt from content (first 200 characters)
		if len(content) > 200 {
			excerpt = content[:200] + "..."
		} else {
			excerpt = content
		}
	}
	if status == "" {
		status = "published"
	}

	// Validate input
	if title == "" || content == "" {
		c.showCreateWithError(w, r, "Title and content are required", title, content)
		return
	}

	// Create blog
	_, err = c.BlogModel.Create(title, content, excerpt, status, user.ID)
	if err != nil {
		c.showCreateWithError(w, r, "Failed to create blog", title, content)
		return
	}

	// Redirect to blogs list
	http.Redirect(w, r, "/dashboard/blogs", http.StatusSeeOther)
}

// Edit shows the edit blog form
func (c *BlogController) Edit(w http.ResponseWriter, r *http.Request) {
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

	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user can edit this blog
	canEdit, err := c.BlogModel.CanUserEdit(id, user.ID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if !canEdit {
		http.Error(w, "You don't have permission to edit this blog", http.StatusForbidden)
		return
	}

	// Get blog
	blog, err := c.BlogModel.GetByID(id)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Title": "Edit Blog",
		"Blog":  blog,
		"User":  user,
	}

	renderTemplate(w, "dashboard/blogs/edit", data)
}

// Update updates an existing blog post
func (c *BlogController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard/blogs", http.StatusSeeOther)
		return
	}

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

	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user can edit this blog
	canEdit, err := c.BlogModel.CanUserEdit(id, user.ID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if !canEdit {
		http.Error(w, "You don't have permission to edit this blog", http.StatusForbidden)
		return
	}

	// Get form data
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	excerpt := strings.TrimSpace(r.FormValue("excerpt"))
	status := strings.TrimSpace(r.FormValue("status"))

	// Set defaults
	if excerpt == "" {
		// Auto-generate excerpt from content (first 200 characters)
		if len(content) > 200 {
			excerpt = content[:200] + "..."
		} else {
			excerpt = content
		}
	}
	if status == "" {
		status = "published"
	}

	// Validate input
	if title == "" || content == "" {
		c.showEditWithError(w, r, id, "Title and content are required", title, content)
		return
	}

	// Update blog
	_, err = c.BlogModel.Update(id, title, content, excerpt, status)
	if err != nil {
		c.showEditWithError(w, r, id, "Failed to update blog", title, content)
		return
	}

	// Redirect to blogs list
	http.Redirect(w, r, "/dashboard/blogs", http.StatusSeeOther)
}

// Delete deletes a blog post
func (c *BlogController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard/blogs", http.StatusSeeOther)
		return
	}

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

	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user can delete this blog (owner or admin)
	canDelete, err := c.BlogModel.CanUserDelete(id, user.ID, user.Role)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if !canDelete {
		http.Error(w, "You don't have permission to delete this blog", http.StatusForbidden)
		return
	}

	// Delete blog
	err = c.BlogModel.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	// Redirect to blogs list
	http.Redirect(w, r, "/dashboard/blogs", http.StatusSeeOther)
}

// AdminDelete deletes any blog post (admin only)
func (c *BlogController) AdminDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard/admin/blogs", http.StatusSeeOther)
		return
	}

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

	// Get current user
	user, err := middleware.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user is admin
	if !user.IsAdmin() {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Delete blog (admin can delete any blog)
	err = c.BlogModel.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	// Redirect to admin blogs list
	http.Redirect(w, r, "/dashboard/admin/blogs", http.StatusSeeOther)
}

// showCreateWithError displays create form with error
func (c *BlogController) showCreateWithError(w http.ResponseWriter, r *http.Request, errorMsg, title, content string) {
	user, _ := middleware.GetCurrentUser(r)
	data := map[string]interface{}{
		"Title":      "Create New Blog",
		"User":       user,
		"Error":      errorMsg,
		"OldTitle":   title,
		"OldContent": content,
	}
	renderTemplate(w, "dashboard/blogs/create", data)
}

// showEditWithError displays edit form with error
func (c *BlogController) showEditWithError(w http.ResponseWriter, r *http.Request, id int, errorMsg, title, content string) {
	user, _ := middleware.GetCurrentUser(r)
	blog := &models.Blog{
		ID:      id,
		Title:   title,
		Content: content,
	}
	data := map[string]interface{}{
		"Title": "Edit Blog",
		"User":  user,
		"Blog":  blog,
		"Error": errorMsg,
	}
	renderTemplate(w, "dashboard/blogs/edit", data)
}
