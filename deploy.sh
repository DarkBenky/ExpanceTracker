#!/bin/bash

# ExpenseTracker Deployment Script
echo "🚀 Starting ExpenseTracker deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    print_error "Please do not run this script as root"
    exit 1
fi

# Get server IP
if [ -z "$1" ]; then
    print_error "Usage: $0 <SERVER_IP>"
    print_error "Example: $0 192.168.1.100"
    exit 1
fi

SERVER_IP=$1
print_status "Deploying to server IP: $SERVER_IP"

# Update .env file with server IP
print_status "Updating .env file with server IP..."
cat > frontend/.env << EOF
NODE_ENV=production
VUE_APP_API_URL=http://$SERVER_IP:1234/
EOF

# Build frontend
print_status "Building frontend for production..."
cd frontend
npm run build
if [ $? -ne 0 ]; then
    print_error "Frontend build failed!"
    exit 1
fi
cd ..

# Build Go backend
print_status "Building Go backend..."
go build -o expensetracker main.go
if [ $? -ne 0 ]; then
    print_error "Go build failed!"
    exit 1
fi

print_status "✅ Build completed successfully!"
print_status "📁 Frontend files are in: frontend/dist/"
print_status "📁 Backend binary is: expensetracker"

echo ""
print_status "🎯 Next steps for deployment:"
echo "1. Copy frontend/dist/ to your web server"
echo "2. Copy expensetracker binary to your server"
echo "3. Run the backend: ./expensetracker"
echo "4. Serve frontend on port 8080 or configure web server"
echo ""
print_status "🔧 Quick deployment commands:"
echo "Backend: ./expensetracker"
echo "Frontend: python3 -m http.server 8080 (from dist/ directory)"
