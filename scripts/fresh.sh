#!/bin/bash
# scripts/fresh.sh - Fresh install (like php artisan migrate:fresh --seed)

echo "🆕 Fresh database setup..."

# Drop and recreate database
DB_NAME=$(grep DB_NAME .env | cut -d'=' -f2)
DB_USER=$(grep DB_USER .env | cut -d'=' -f2)
DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d'=' -f2)
DB_HOST=$(grep DB_HOST .env | cut -d'=' -f2)

echo "🗑️  Dropping existing database..."
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD -e "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null

echo "🔧 Creating fresh database..."
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD -e "CREATE DATABASE $DB_NAME;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ Database recreated successfully"
else
    echo "❌ Database creation failed"
    exit 1
fi

# Run migrations
echo "🔄 Running migrations..."
./scripts/migrate.sh

# Run seeders
echo "🌱 Seeding database..."
./scripts/seed.sh

echo "🎉 Fresh database setup completed!"
