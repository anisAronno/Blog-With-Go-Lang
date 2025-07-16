package middleware

import (
	"context"
	"go-web-app/app/models"
	"go-web-app/config"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

// SessionStore holds the session store
var SessionStore *sessions.CookieStore

// InitSessions initializes the session store
func InitSessions() {
	SessionStore = sessions.NewCookieStore([]byte(config.AppConfig.SessionSecret))
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}
}

// AuthMiddleware checks if user is authenticated
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := SessionStore.Get(r, "session")
		if err != nil {
			log.Printf("Session error: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user ID to context for use in handlers
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GuestMiddleware redirects authenticated users away from guest pages
func GuestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := SessionStore.Get(r, "session")
		if err != nil {
			// If session error, continue as guest
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := session.Values["user_id"]
		if ok && userID != nil {
			// User is authenticated, redirect to dashboard
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(r *http.Request) (*models.User, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return nil, nil
	}

	// Convert userID to int
	var id int
	switch v := userID.(type) {
	case int:
		id = v
	case string:
		var err error
		id, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	// Get user from database
	userModel := models.NewUserModel(config.Database)
	return userModel.GetByID(id)
}

// GetCurrentUserFromSession returns user from session (for non-protected routes)
func GetCurrentUserFromSession(r *http.Request) (*models.User, error) {
	session, err := SessionStore.Get(r, "session")
	if err != nil {
		return nil, nil // Don't error on homepage
	}

	userID, ok := session.Values["user_id"]
	if !ok || userID == nil {
		return nil, nil
	}

	// Convert userID to int
	var id int
	switch v := userID.(type) {
	case int:
		id = v
	case string:
		var err error
		id, err = strconv.Atoi(v)
		if err != nil {
			return nil, nil
		}
	default:
		return nil, nil
	}

	// Get user from database
	userModel := models.NewUserModel(config.Database)
	return userModel.GetByID(id)
}

// SetUserSession sets user session data
func SetUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := SessionStore.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = userID
	return session.Save(r, w)
}

// ClearUserSession clears user session data
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := SessionStore.Get(r, "session")
	if err != nil {
		return err
	}

	// Clear session values
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1 // Mark for deletion

	return session.Save(r, w)
}
