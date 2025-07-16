#!/bin/bash
# scripts/seed.sh - Run database seeders (like php artisan db:seed)

echo "🌱 Seeding database with sample data..."

# Check if migrations have been run
DB_NAME=$(grep DB_NAME .env | cut -d'=' -f2)
DB_USER=$(grep DB_USER .env | cut -d'=' -f2)
DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d'=' -f2)
DB_HOST=$(grep DB_HOST .env | cut -d'=' -f2)

echo "📍 Database: $DB_NAME"

# Check if tables exist
TABLES=$(mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME -e "SHOW TABLES;" 2>/dev/null | wc -l)

if [ $TABLES -lt 2 ]; then
    echo "⚠️  No tables found. Running migrations first..."
    ./scripts/migrate.sh
fi

# Run seeders
echo "🔄 Running seeders..."
go run database/seeders/seed.go

if [ $? -eq 0 ]; then
    echo "✅ Database seeded successfully!"
    echo ""
    echo "🎯 Demo user accounts created:"
    echo "   - Email: admin@example.com | Password: password"
    echo "   - Email: john@example.com  | Password: password"
    echo "   - Email: jane@example.com  | Password: password"
else
    echo "❌ Seeding failed"
    exit 1
fi
