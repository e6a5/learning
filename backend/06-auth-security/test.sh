#!/bin/bash

echo "🔐 Testing Authentication & Security Server"
echo "============================================"

echo ""
echo "1. Testing server status..."
curl -s http://localhost:8081/ | jq .

echo ""
echo "2. Testing admin login..."
curl -s -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq .

echo ""
echo "3. Testing user registration..."
curl -s -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"securepass123","email":"test@example.com"}' | jq .

echo ""
echo "4. Testing new user login..."
curl -s -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"securepass123"}' | jq . 