# 🔴 03 - Redis Introduction

This project demonstrates how to connect Go with Redis and build a REST API for caching and key-value operations. Learn the fundamentals of Redis through practical examples.

---

## 🎯 What You'll Learn

- Connect Go to Redis using the `go-redis` client
- Implement basic Redis operations (GET, SET, DELETE)
- Work with TTL (Time To Live) for cache expiration
- Use Docker Compose to run Redis locally
- Build a fully containerized application

---

## 🧱 Stack

- **Golang** - Backend server
- **Redis** - In-memory data store and cache
- **Docker & Docker Compose** - Containerization
- **Gorilla Mux** - HTTP router

---

## 🚀 Quick Start

```bash
# Start Redis and the API
make up

# Check if everything is running
make ps

# Test the health endpoint
make test-health
```

## 🔄 API Endpoints

| Method | Endpoint              | Description                    | Example                              |
| ------ | --------------------- | ------------------------------ | ------------------------------------ |
| GET    | `/health`             | Health check                   | `curl /health`                       |
| POST   | `/cache`              | Set a key-value pair           | `{"key":"name","value":"John","ttl":300}` |
| GET    | `/cache/{key}`        | Get value by key               | `curl /cache/name`                   |
| DELETE | `/cache/{key}`        | Delete a key                   | `curl -X DELETE /cache/name`         |
| GET    | `/cache`              | Get all keys                   | `curl /cache`                        |
| GET    | `/cache/{key}/ttl`    | Get TTL for a key              | `curl /cache/name/ttl`               |
| POST   | `/cache/{key}/expire` | Set expiration for existing key | `{"ttl":600}`                       |

---

## 🧪 Test It Out

```bash
# 1. Set a key with TTL
make test-set

# 2. Get the key
make test-get

# 3. Check TTL
make test-ttl

# 4. List all keys
make test-keys

# 5. Delete the key
make test-delete
```

### Manual Testing Examples

```bash
# Set a session token with 1 hour expiration
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{"key":"session:user123","value":"abc123def","ttl":3600}'

# Get user preferences (permanent storage)
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{"key":"user:123:prefs","value":"{\"theme\":\"dark\",\"lang\":\"en\"}"}'

# Check if key exists and get its TTL
curl http://localhost:8080/cache/session:user123/ttl

# Set expiration on an existing key (extend session)
curl -X POST http://localhost:8080/cache/session:user123/expire \
  -H "Content-Type: application/json" \
  -d '{"ttl":7200}'
```

---

## 🔧 Development Commands

```bash
# Build and start services
make up

# View application logs
make logs

# View Redis logs
make redis-logs

# Access Redis CLI directly
make redis-cli

# Stop services
make down

# Clean everything and rebuild
make rebuild

# See all available commands
make help
```

---

## 🔍 Redis Operations Explained

### Key-Value Storage
- **SET**: Store string values with optional expiration
- **GET**: Retrieve values by key
- **DEL**: Remove keys from Redis

### TTL (Time To Live)
- Set expiration when creating keys
- Check remaining time with TTL command
- Extend expiration for existing keys

### Use Cases
- **Session Storage**: User sessions with automatic cleanup
- **Caching**: Store frequently accessed data
- **Rate Limiting**: Track API usage with expiring counters
- **Temporary Data**: Store data that should auto-expire

---

## 📦 File Overview

| File            | Purpose                              |
| --------------- | ------------------------------------ |
| `main.go`       | Go server with Redis operations      |
| `go.mod`        | Go module and dependencies           |
| `Dockerfile`    | Multi-stage build for the Go app     |
| `compose.yml`   | Docker Compose for Redis + App      |
| `Makefile`      | Development and testing commands     |

---

## 🐳 Docker Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Go App        │────│   Redis         │
│   :8080         │    │   :6379         │
│                 │    │   (persistent)  │
└─────────────────┘    └─────────────────┘
```

- **App Container**: Runs the Go application
- **Redis Container**: Persistent Redis with AOF enabled
- **Network**: Isolated Docker network for service communication
- **Volume**: Persistent storage for Redis data

---

## 💡 Next Steps

Try these experiments:
1. **Performance Testing**: Use Apache Bench to test Redis vs database speed
2. **Advanced Data Types**: Explore Redis lists, sets, and hashes
3. **Pub/Sub**: Implement real-time messaging with Redis
4. **Clustering**: Set up Redis in cluster mode for high availability 