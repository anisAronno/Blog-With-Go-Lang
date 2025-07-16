#!/bin/bash
# scripts/migrate.sh - Run database migrations (like php artisan migrate)

echo "🗄️  Running database migrations..."

# Check if database exists
DB_NAME=$(grep DB_NAME .env | cut -d'=' -f2)
DB_USER=$(grep DB_USER .env | cut -d'=' -f2)
DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d'=' -f2)
DB_HOST=$(grep DB_HOST .env | cut -d'=' -f2)

echo "📍 Database: $DB_NAME"
echo "🏠 Host: $DB_HOST"

# Create database if it doesn't exist
echo "🔧 Creating database if it doesn't exist..."
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ Database connection successful"
else
    echo "❌ Database connection failed. Please check your .env configuration"
    exit 1
fi

# Run migrations
echo "🔄 Running migrations..."
go run database/migrations/migrate.go

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed successfully!"
else
    echo "❌ Migration failed"
    exit 1
fi
