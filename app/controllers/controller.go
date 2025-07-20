// app/controllers/controller.go - Base controller with common functionality
package controllers

import (
	"html/template"
	"net/http"
	"strings"
)

// renderTemplate renders a template with the given data
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Determine layout and template execution based on template path
	var layoutPath string
	var executeTemplate string

	if strings.HasPrefix(tmpl, "dashboard/") {
		// Dashboard templates use dashboard layout
		layoutPath = "templates/dashboard/layout.html"
		executeTemplate = "dashboard_layout"
	} else {
		// Other templates use base layout
		layoutPath = "templates/layouts/base.html"
		executeTemplate = "base"
	}

	templatePath := "templates/" + tmpl + ".html"

	// Component template paths
	componentPaths := []string{
		"templates/components/header.html",
		"templates/components/dashboard-header.html",
		"templates/components/footer.html",
		"templates/components/pagination.html",
	}

	// Create template with helper functions
	t := template.New(executeTemplate).Funcs(template.FuncMap{
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

	// Prepare all template files
	allTemplateFiles := []string{layoutPath, templatePath}
	allTemplateFiles = append(allTemplateFiles, componentPaths...)

	// Parse template files
	t, err := t.ParseFiles(allTemplateFiles...)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the appropriate template
	err = t.ExecuteTemplate(w, executeTemplate, data)
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
