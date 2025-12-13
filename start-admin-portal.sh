#!/bin/bash

echo "🚀 Starting Collect Your World Admin Portal"
echo "========================================="

# Navigate to project root
cd /Users/ntho/Projects/Collect-your-world/service-platform

echo "📦 Starting infrastructure (PostgreSQL, Redis)..."
task up-infra

echo "⏳ Waiting for services to start..."
sleep 5

echo "🗄️ Running database migrations..."
task clean-migrate

echo "🔧 Starting backend API server..."
task run &
BACKEND_PID=$!

echo "⏳ Waiting for backend to start..."
sleep 3

# Test if backend is running
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Backend is running successfully"
else
    echo "❌ Backend failed to start. Please check logs."
    exit 1
fi

echo "🌐 Starting admin portal frontend..."
cd admin-portal

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing npm dependencies..."
    npm install
fi

npm start &
FRONTEND_PID=$!

echo ""
echo "🎉 All services started successfully!"
echo "========================================="
echo "📍 Admin Portal: http://localhost:3000"
echo "📍 Backend API:  http://localhost:8080"
echo "📍 API Docs:     http://localhost:8080/api/v1/swagger/*"
echo ""
echo "🔧 Backend PID: $BACKEND_PID"
echo "🔧 Frontend PID: $FRONTEND_PID"
echo ""
echo "Press Ctrl+C to stop all services"
echo "========================================="

# Wait for interrupt
trap 'echo "🛑 Stopping services..."; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit' INT

# Keep script running
wait