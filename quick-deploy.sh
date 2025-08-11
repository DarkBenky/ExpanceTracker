#!/bin/bash

# Quick deployment script
SERVER_IP=${1:-"localhost"}

echo "🚀 Deploying ExpenseTracker to $SERVER_IP"

# Update frontend config
echo "NODE_ENV=production
VUE_APP_API_URL=http://$SERVER_IP:1234/" > frontend/.env

# Build frontend
echo "📦 Building frontend..."
cd frontend && npm run build && cd ..

# Build backend
echo "🔨 Building backend..."
go build -o expensetracker main.go

echo "✅ Build complete!"
echo ""
echo "📋 To deploy:"
echo "1. Copy frontend/dist/ to your web server"
echo "2. Copy expensetracker binary to your server"
echo "3. Run: ./expensetracker (backend)"
echo "4. Serve frontend/dist/ on port 8080"
echo ""
echo "🌐 Access at: http://$SERVER_IP:8080"
