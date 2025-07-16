#!/bin/bash
# scripts/test.sh - Run test suite (like php artisan test)

echo "🧪 Running Go Web App Test Suite..."
echo "==================================="

# Check if test files exist
if [ ! -d "tests" ]; then
    echo "❌ Tests directory not found"
    exit 1
fi

# Run tests with verbose output
echo "🔄 Running tests..."
go test ./tests/... -v

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All tests passed successfully!"
else
    echo ""
    echo "❌ Some tests failed"
    exit 1
fi

# Run benchmarks if available
echo ""
echo "📊 Running benchmarks..."
go test ./tests/... -bench=. -benchmem

echo ""
echo "🎯 Test coverage report:"
go test ./tests/... -cover
