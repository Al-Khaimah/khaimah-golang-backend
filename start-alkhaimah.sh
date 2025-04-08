#!/bin/bash

echo "🚀 Starting Al-Khaimah Application..."

echo "📦 Starting PostgresSQL Docker services..."
docker-compose up -d

echo "📦 Installing Go packages (if needed)..."
go mod tidy

echo "🔨 Building the application..."
go build -o alkhaimah cmd/*.go

echo "🏃 Running the application..."
nohup ./alkhaimah > alkhaimah.log 2>&1 &

echo "✅ Al-Khaimah is running in the background. Check logs with: tail -f alkhaimah.log"
