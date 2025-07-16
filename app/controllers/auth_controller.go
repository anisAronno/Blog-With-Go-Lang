package controllers

import (
	"go-web-app/app/middleware"
	"go-web-app/app/models"
	"go-web-app/config"
	"net/http"
	"strings"
)

// AuthController handles authentication related requests
type AuthController struct {
	UserModel *models.UserModel
}

// NewAuthController creates a new AuthController
func NewAuthController() *AuthController {
	return &AuthController{
		UserModel: models.NewUserModel(config.Database),
	}
}

// ShowLogin displays the login form
func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login",
	}

	renderTemplate(w, "auth/login", data)
}

// ShowRegister displays the registration form
func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Register",
	}

	renderTemplate(w, "auth/register", data)
}

// Login handles user login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		c.showLoginWithError(w, r, "Email and password are required")
		return
	}

	// Authenticate user
	user, err := c.UserModel.Authenticate(email, password)
	if err != nil {
		c.showLoginWithError(w, r, "Invalid email or password")
		return
	}

	// Set session
	err = middleware.SetUserSession(w, r, user.ID)
	if err != nil {
		c.showLoginWithError(w, r, "Failed to create session")
		return
	}

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Register handles user registration
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password_confirmation")

	// Validate input
	if name == "" || email == "" || password == "" {
		c.showRegisterWithError(w, r, "All fields are required")
		return
	}

	if password != passwordConfirmation {
		c.showRegisterWithError(w, r, "Passwords do not match")
		return
	}

	if len(password) < 6 {
		c.showRegisterWithError(w, r, "Password must be at least 6 characters long")
		return
	}

	// Check if email already exists
	exists, err := c.UserModel.EmailExists(email)
	if err != nil {
		c.showRegisterWithError(w, r, "Failed to check email availability")
		return
	}

	if exists {
		c.showRegisterWithError(w, r, "Email already registered")
		return
	}

	// Create user
	user, err := c.UserModel.Create(name, email, password)
	if err != nil {
		c.showRegisterWithError(w, r, "Failed to create account")
		return
	}

	// Set session
	err = middleware.SetUserSession(w, r, user.ID)
	if err != nil {
		c.showRegisterWithError(w, r, "Account created but failed to login")
		return
	}

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout handles user logout
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	err := middleware.ClearUserSession(w, r)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// showLoginWithError displays login form with error message
func (c *AuthController) showLoginWithError(w http.ResponseWriter, r *http.Request, errorMsg string) {
	data := map[string]interface{}{
		"Title": "Login",
		"Error": errorMsg,
		"Email": r.FormValue("email"), // Preserve email input
	}

	renderTemplate(w, "auth/login", data)
}

// showRegisterWithError displays register form with error message
func (c *AuthController) showRegisterWithError(w http.ResponseWriter, r *http.Request, errorMsg string) {
	data := map[string]interface{}{
		"Title": "Register",
		"Error": errorMsg,
		"Name":  r.FormValue("name"),  // Preserve form data
		"Email": r.FormValue("email"),
	}

	renderTemplate(w, "auth/register", data)
}
