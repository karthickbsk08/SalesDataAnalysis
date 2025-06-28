#!/bin/bash

set -e

echo "🚀 Starting Golang WebAuthn Project Setup..."

# Step 1: Check Go installation
if ! command -v go &> /dev/null
then
    echo "❌ Go is not installed. Please install Go 1.20+ and try again."
    exit 1
fi

# Step 2: Clone the repository (skip if already cloned)
REPO_URL="https://github.com/karthickbsk08/SalesDataAnalysis.git"

REPO_NAME=$(basename -s .git "$REPO_URL")
echo "Repository name: $REPO_NAME"

WORKDIR="$(pwd)/$REPO_NAME"
echo "Current working directory: $WORKDIR"


# Step 3: check if working directory exist in current workspace, if not create, else move into directory and clone the URL.
if [ ! -d "$WORKDIR" ]; then
echo "📂 Creating working directory at $WORKDIR (if not exists)..."
mkdir -p "$WORKDIR"
cd "$WORKDIR"
echo "📁 Now inside $WORKDIR"
    echo "📥 Cloning repository..."
    git clone "$REPO_URL"
fi

cd "$WORKDIR"

#step 5 : Initialize and mod 
#echo  "create a go mod file"
#go mod init "$PROJECTNAME"

# Step 4: Initialize and install dependencies
echo "📦 Running go mod tidy..."
go mod tidy 

# Step 4: Setup environment
if [ ! -f ".env" ]; then
    echo "⚙️ Creating .env file from example..."
        cp .env.example .env
fi

# Step 5: Run the application
echo "🚀 Starting the application..."
go run main.go

