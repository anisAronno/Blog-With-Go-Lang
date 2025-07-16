# Go Web App Development Guide

## Quick Start

### 1. Development Setup

```bash
# Install dependencies and development tools
make install

# Setup database (run migrations and seeders)
make db-setup

# Start development server with live reload
make dev
```

### 2. Alternative: Docker Development

```bash
# Start development environment with Docker
make docker-dev

# Stop development environment
make docker-dev-down
```

## Available Commands

### Database Commands

```bash
make migrate            # Run all pending migrations
make migrate-rollback   # Rollback last migration
make migrate-status     # Show migration status
make migrate-reset      # Reset all migrations

make seed              # Run all seeders
make seed-users        # Run only user seeder
make seed-blogs        # Run only blog seeder

make db-setup          # Run migrations then seeders
make db-reset          # Reset database completely
make fresh             # Fresh database setup
```

### Development Commands

```bash
make dev               # Start with live reload (Air)
make serve             # Start without live reload
make test              # Run all tests
make build             # Build for production
make clean             # Clean build files
```

### Docker Commands

```bash
make docker-dev        # Development with live reload
make docker-up         # Production Docker setup
make docker-down       # Stop Docker containers
make docker-build      # Build Docker images
make docker-clean      # Clean Docker artifacts
```

## Features

### 1. ✅ Role-Based Access Control

- **Admin**: Can manage all users, blogs, and system settings
- **Author**: Can create and manage their own blogs
- **User**: Can view content (basic user)

### 2. ✅ Protected Admin Account

- Main admin account (`admin@example.com`) cannot be deleted
- Admins cannot delete their own accounts
- Additional safety checks in place

### 3. ✅ Modern Dashboard UI

- Responsive design with Tailwind CSS
- Common navigation header for all dashboard pages
- Mobile-friendly with dropdown menus
- Clean and professional interface

### 4. ✅ Database Management

- Modular migration system with versioning
- Separate seeder classes for different data types
- CLI tools for easy database management
- Make commands for quick operations

### 5. ✅ Development Environment

- Live reload with Air
- Docker development setup
- Volume mounting for instant code changes
- Development-specific configurations

## Project Structure

```
├── app/
│   ├── controllers/     # HTTP request handlers
│   ├── middleware/      # Authentication & other middleware
│   ├── models/         # Database models
│   └── tests/          # Application tests
├── cmd/
│   ├── migrate/        # Migration CLI tool
│   └── seed/           # Seeder CLI tool
├── database/
│   ├── migrations/     # Database migrations
│   └── seeders/        # Database seeders
├── templates/
│   ├── dashboard/      # Dashboard templates with common layout
│   └── ...             # Other templates
├── docker-compose.yml          # Production Docker setup
├── docker-compose.dev.yml      # Development overrides
├── Dockerfile                  # Production Dockerfile
├── Dockerfile.dev             # Development Dockerfile
├── .air.toml                  # Live reload configuration
└── Makefile                   # Development commands
```

## Environment Variables

Create a `.env` file with:

```env
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=go_web_app
DB_USERNAME=root
DB_PASSWORD=your_password
SESSION_SECRET=your-secret-key
APP_ENV=development
```

## Default Users (after seeding)

- **Admin**: admin@example.com / admin123
- **Author**: john@example.com / password123
- **User**: jane@example.com / password123

## Live Development

The development setup supports live reloading, so any code changes will automatically restart the server and reflect in the browser immediately.

```bash
# Start development with live reload
make dev

# Or with Docker for consistent environment
make docker-dev
```

Visit http://localhost:3000 to see your changes in real-time!
