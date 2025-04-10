#!/bin/bash

echo "🚀 Starting Al-Khaimah Application..."

# 1. Stash any local changes
echo "📦 Stashing local changes..."
git stash --include-untracked

# 2. Pull the latest changes from GitHub
echo "📥 Pulling latest changes from GitHub..."
git pull origin master

# 3. Start PostgreSQL via Docker (if any)
echo "🐘 Starting PostgreSQL Docker service (if any)..."
docker-compose up -d

# 4. Tidy Go modules
echo "📦 Tidying Go modules..."
go mod tidy

# 5. Stop existing Go server
echo "🛑 Stopping existing server (if running)..."
pkill sudo -f alkhaimah || true

# 6. Build the application
echo "🔨 Building the application..."
go build -o alkhaimah cmd/*.go

# 7. Start the server in the background
echo "🏃 Starting the application in the background..."
nohup ./alkhaimah > alkhaimah.log 2>&1 &

# 8. Try applying the stash safely
if git stash list | grep -q .; then
  echo "🎯 Trying to reapply stash..."

  if git stash pop --quiet; then
    echo "✅ Stash applied successfully!"
  else
    echo "❌ Conflict detected while applying stash. Dropping it..."
    git reset --hard
    git stash drop
  fi
else
  echo "🎯 No stash to apply."
fi

# 9. Done!
echo "✅ Al-Khaimah is running. View logs with: tail -f alkhaimah.log"
