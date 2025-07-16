#!/bin/bash
# scripts/setup.sh - Initial project setup (like Laravel installation)

echo "🔧 Setting up Go Web App..."
echo "=================================="

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go $(go version | cut -d' ' -f3) found"

# Check MySQL installation
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL client is not installed. Please install MySQL first."
    exit 1
fi

echo "✅ MySQL client found"

# Install dependencies
echo ""
echo "1️⃣  Installing Go dependencies..."
go mod tidy
go mod download

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully!"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Create database
echo ""
echo "2️⃣  Creating database..."
DB_NAME=$(grep DB_NAME .env | cut -d'=' -f2)
DB_USER=$(grep DB_USER .env | cut -d'=' -f2)
DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d'=' -f2)
DB_HOST=$(grep DB_HOST .env | cut -d'=' -f2)

mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ Database created successfully!"
else
    echo "❌ Database creation failed. Please check your MySQL credentials in .env file"
    echo "💡 Make sure MySQL server is running and credentials are correct"
    exit 1
fi

# Run migrations
echo ""
echo "3️⃣  Running database migrations..."
./scripts/migrate.sh

# Seed database
echo ""
echo "4️⃣  Seeding database with sample data..."
./scripts/seed.sh

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "🚀 To start the development server, run:"
echo "   make serve"
echo "   OR"
echo "   ./scripts/serve.sh"
echo ""
echo "📝 Available commands:"
echo "   make help              - Show all available commands"
echo "   make serve             - Start development server"
echo "   make migrate           - Run migrations"
echo "   make seed              - Seed database"
echo "   make fresh             - Fresh database setup"
echo "   make test              - Run tests"
echo ""
echo "🌐 Once server is running, visit: http://localhost:3000"
