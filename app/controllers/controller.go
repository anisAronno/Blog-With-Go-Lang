// app/controllers/controller.go - Base controller with common functionality
package controllers

import (
	"html/template"
	"net/http"
	"strings"
)

// renderTemplate renders a template with the given data
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse base layout and specific template
	layoutPath := "templates/layouts/base.html"
	templatePath := "templates/" + tmpl + ".html"

	// Create template with helper functions
	t := template.New("base").Funcs(template.FuncMap{
		"split": strings.Split,
		"slice": func(s string, start, end int) string {
			if start >= len(s) {
				return ""
			}
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"add": func(a, b int) int {
			return a + b
		},
		"asset": func(path string) string {
			return "/public/" + path
		},
	})

	// Parse template files
	t, err := t.ParseFiles(layoutPath, templatePath)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
	}
}

// StaticFileHandler serves static files from the public directory
func StaticFileHandler() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("./public/")))
}

// Helper function to format templates with context
func formatTemplate(name string, data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Add common template functions
	data["asset"] = func(path string) string {
		return "/public/" + path
	}

	return data
}
