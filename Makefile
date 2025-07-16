# Go Web App - Laravel-style commands
# Make commands for easy project management

.PHONY: help install migrate seed serve test clean build docker-up docker-down fresh

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies (like composer install)
	@echo "📦 Installing Go dependencies..."
	go mod tidy
	go mod download
	@echo "🔧 Installing development tools..."
	@go install github.com/cosmtrek/air@latest 2>/dev/null || echo "⚠️  Air installation failed, live reload may not work"
	@echo "✅ Dependencies installed successfully!"

migrate: ## Run database migrations (like php artisan migrate)
	@echo "🗄️  Running database migrations..."
	go run cmd/migrate/main.go up
	@echo "✅ Migrations completed successfully!"

migrate-rollback: ## Rollback last migration
	@echo "⬇️  Rolling back last migration..."
	go run cmd/migrate/main.go down
	@echo "✅ Migration rollback completed!"

migrate-status: ## Show migration status
	@echo "📊 Checking migration status..."
	go run cmd/migrate/main.go status

migrate-reset: ## Reset all migrations (rollback all then migrate)
	@echo "🔄 Resetting all migrations..."
	go run cmd/migrate/main.go down
	go run cmd/migrate/main.go up
	@echo "✅ Migration reset completed!"

seed: ## Run database seeders (like php artisan db:seed)
	@echo "🌱 Seeding database with sample data..."
	go run cmd/seed/main.go
	@echo "✅ Database seeded successfully!"

seed-users: ## Run only user seeder
	@echo "👥 Seeding users..."
	go run cmd/seed/main.go users
	@echo "✅ Users seeded successfully!"

seed-blogs: ## Run only blog seeder
	@echo "📝 Seeding blogs..."
	go run cmd/seed/main.go blogs
	@echo "✅ Blogs seeded successfully!"

serve: ## Start the development server (like php artisan serve)
	@echo "🚀 Starting development server..."
	@echo "📝 Server will be available at: http://localhost:3000"
	@echo "🛑 Press Ctrl+C to stop the server"
	go run main.go

test: ## Run tests (like php artisan test)
	@echo "🧪 Running test suite..."
	go test ./tests/... -v
	@echo "✅ Tests completed!"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	go clean
	rm -f go-web-app
	@echo "✅ Cleanup completed!"

build: ## Build the application for production
	@echo "🔨 Building application..."
	go build -o go-web-app main.go
	@echo "✅ Build completed! Executable: ./go-web-app"

fresh: ## Fresh install (like php artisan migrate:fresh --seed)
	@echo "🆕 Fresh database setup..."
	make migrate-reset
	make seed
	@echo "✅ Fresh database setup completed!"

db-setup: ## Setup database (run migrations then seeders)
	@echo "🗄️  Setting up database..."
	make migrate
	make seed
	@echo "✅ Database setup completed!"

db-reset: ## Reset database (reset migrations then seed)
	@echo "🔄 Resetting database..."
	make migrate-reset
	make seed
	@echo "✅ Database reset completed!"

setup: ## Initial project setup
	@echo "🔧 Setting up Go Web App..."
	@echo "1️⃣  Installing dependencies..."
	make install
	@echo "2️⃣  Creating database..."
	mysql -u root -pbs@123 -e "CREATE DATABASE IF NOT EXISTS go_web_app;" || echo "Database already exists or connection failed"
	@echo "3️⃣  Running migrations..."
	make migrate
	@echo "4️⃣  Seeding database..."
	make seed
	@echo "🎉 Setup completed! Run 'make serve' to start the server"

db-create: ## Create database
	@echo "🗄️  Creating database..."
	mysql -u root -pbs@123 -e "CREATE DATABASE IF NOT EXISTS go_web_app;"
	@echo "✅ Database created successfully!"

db-drop: ## Drop database
	@echo "🗑️  Dropping database..."
	mysql -u root -pbs@123 -e "DROP DATABASE IF EXISTS go_web_app;"
	@echo "✅ Database dropped successfully!"

# Docker commands (if using Docker)
docker-up: ## Start Docker services
	@if [ -f docker-compose.yml ]; then \
		echo "🐳 Starting Docker services..."; \
		docker-compose up -d; \
	else \
		echo "❌ docker-compose.yml not found"; \
	fi

docker-down: ## Stop Docker services
	@if [ -f docker-compose.yml ]; then \
		echo "🐳 Stopping Docker services..."; \
		docker-compose down; \
	else \
		echo "❌ docker-compose.yml not found"; \
	fi

docker-dev: ## Start development environment with live reload
	@echo "🐳 Starting development environment with live reload..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

docker-dev-down: ## Stop development environment
	@echo "🛑 Stopping development environment..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

docker-build: ## Build Docker images
	@echo "🔨 Building Docker images..."
	docker-compose build

docker-clean: ## Clean Docker containers and images
	@echo "🧹 Cleaning Docker containers and images..."
	docker-compose down --rmi all --volumes --remove-orphans

# Quick development workflow
dev: ## Start development server with live reload (using Air)
	@echo "🚀 Starting development server with live reload..."
	@echo "📝 Server will be available at: http://localhost:3000"
	@echo "🔄 Code changes will trigger automatic reload"
	@echo "🛑 Press Ctrl+C to stop the server"
	@command -v air >/dev/null 2>&1 || { echo "❌ Air not installed. Install with: go install github.com/cosmtrek/air@latest"; exit 1; }
	air -c .air.toml

dev-setup: ## Quick development setup (install + migrate + seed + serve)
	@echo "⚡ Quick development setup..."
	make install
	make fresh
	@echo "🚀 Starting development server..."
	make serve

# Production deployment
deploy: ## Build and prepare for production
	@echo "🚀 Preparing for production deployment..."
	make clean
	make install
	make build
	@echo "✅ Production build ready!"

# Show environment information
env: ## Show environment information
	@echo "🔧 Environment Information:"
	@echo "Go version: $(shell go version)"
	@echo "OS: $(shell uname -s)"
	@echo "Architecture: $(shell uname -m)"
	@echo "Current directory: $(shell pwd)"
	@echo "Database host: $(shell grep DB_HOST .env | cut -d'=' -f2)"
	@echo "App port: $(shell grep APP_PORT .env | cut -d'=' -f2)"

# Lint and format code
lint: ## Lint and format Go code
	@echo "🧹 Formatting and linting code..."
	go fmt ./...
	go vet ./...
	@echo "✅ Code formatting completed!"
