#!/bin/bash

# Clean up any existing installation
echo "Cleaning up existing installation..."
rm -rf node_modules package-lock.json

# Install dependencies
echo "Installing dependencies..."
npm install

# Create a development build
echo "Creating development build..."
npm run build

# Start the development server
echo "Starting development server..."
npm start 