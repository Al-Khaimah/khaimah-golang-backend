#!/bin/bash

echo "🚀 Starting Al-Khaimah Application..."

# Helper: Check if .env contains IP_ADDRESS (server mode)
IS_SERVER=$(grep -E '^IP_ADDRESS=' .env | awk -F= '{print $2}' | xargs)

if [[ -n "$IS_SERVER" ]]; then
  # 🖥 Server mode (Git operations only)
  echo "📦 Stashing local changes..."
  git stash --include-untracked

  echo "📥 Pulling latest changes from GitHub..."
  git pull origin master
else
  # 💻 Local development mode (Docker services only)
  echo "🐘 Starting PostgreSQL Docker service..."
  docker-compose up -d
fi

# 📦 Common Steps (Both Server & Local)
echo "📦 Tidying Go modules..."
go mod tidy

# 🛑 Stop existing server if running
PID=$(pgrep -f ./alkhaimah)
if [[ -n "$PID" ]]; then
  echo "🛑 Stopping existing server (PID: $PID)..."
  kill "$PID" > /dev/null 2>&1 || true
  sleep 1
fi

# 🧹 Clean up old binary (optional but safe)
echo "🧹 Removing old binary..."
rm -f alkhaimah

# 🔨 Rebuild the binary
echo "🔨 Building the application..."
go build -o alkhaimah cmd/*.go

# 🏃 Run in background (new binary)
echo "🏃 Starting the application in the background..."
setsid ./alkhaimah > alkhaimah.log 2>&1 &

# 🎯 Reapply stash (only on server)
if [[ -n "$IS_SERVER" ]] && git stash list | grep -q .; then
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

# ✅ Final message
echo "✅ Al-Khaimah is running. View logs with: tail -f alkhaimah.log"

# 🪵 Follow logs
tail -f alkhaimah.log
