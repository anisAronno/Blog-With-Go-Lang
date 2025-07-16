#!/bin/bash
# scripts/serve.sh - Start development server (like php artisan serve)

echo "🚀 Starting Go Web App Development Server..."
echo "📍 Server will be available at: http://localhost:3000"
echo ""
echo "📝 Available routes:"
echo "   - GET  /                 (Homepage - Blog listing)"
echo "   - GET  /login            (Login page)"
echo "   - GET  /register         (Register page)"
echo "   - GET  /dashboard        (Dashboard - requires auth)"
echo "   - GET  /dashboard/blogs  (Blog management)"
echo "   - GET  /dashboard/users  (User listing)"
echo "   - GET  /dashboard/profile (User profile)"
echo ""
echo "🛑 Press Ctrl+C to stop the server"
echo "==============================================="

# Run the application
go run main.go
