# 🗃️ 02 - MySQL CRUD API

This project demonstrates how to connect Go with MySQL and build a REST API with full CRUD operations on a `users` table. Learn database fundamentals through practical examples.

---

## 🎯 What You'll Learn

- Connect Go to MySQL using the `database/sql` package
- Write SQL queries (INSERT, SELECT, UPDATE, DELETE)
- Handle HTTP routing with `gorilla/mux`
- Use Docker Compose to run MySQL locally
- Build a fully containerized application
- Implement proper database initialization

---

## 🧱 Stack

- **Golang** - Backend server
- **MySQL 8** - Relational database
- **Docker & Docker Compose** - Containerization
- **Gorilla Mux** - HTTP router

---

## 🚀 Quick Start

```bash
# Start MySQL and the API
make up

# Check if everything is running
make ps

# Test the API
curl http://localhost:8080/users
```

## 🔄 API Endpoints

| Method | Endpoint      | Description       | Example Request                          |
| ------ | ------------- | ----------------- | ---------------------------------------- |
| GET    | `/users`      | List all users    | `curl /users`                           |
| POST   | `/users`      | Create a new user | `{"name":"John","email":"john@test.com"}` |
| PUT    | `/users/{id}` | Update a user     | `{"name":"Jane","email":"jane@test.com"}` |
| DELETE | `/users/{id}` | Delete a user     | `curl -X DELETE /users/1`               |

---

## 🧪 Test It Out

### Quick Tests
```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'

# Get all users
curl http://localhost:8080/users

# Update user with ID 1
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","email":"alice.smith@example.com"}'

# Delete user with ID 1
curl -X DELETE http://localhost:8080/users/1
```

### Complete User Workflow
```bash
# 1. Create multiple users
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","email":"bob@test.com"}'

curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Carol","email":"carol@test.com"}'

# 2. List all users
curl http://localhost:8080/users

# 3. Update a user
curl -X PUT http://localhost:8080/users/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"Carol Johnson","email":"carol.johnson@test.com"}'

# 4. Clean up - delete users
curl -X DELETE http://localhost:8080/users/1
curl -X DELETE http://localhost:8080/users/2
```

---

## 🔧 Development Commands

```bash
# Build and start services
make up

# View application logs
make logs

# View MySQL logs
make db-logs

# Access MySQL CLI directly
make db-cli

# Check running services
make ps

# Stop services
make down

# Clean everything and rebuild
make rebuild

# Quick API test
make curl
```

---

## 🗄️ Database Details

### Schema
```sql
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);
```

### Configuration
- **Database**: `learninglab`
- **User**: `user`
- **Password**: `pass`
- **Port**: `3306`
- **Container**: `mysql-lab`

### Accessing MySQL CLI
```bash
# Via Makefile
make db-cli

# Direct command
docker exec -it mysql-lab mysql -uuser -ppass learninglab
```

---

## 🔍 Implementation Details

### Database Connection
- Uses `database/sql` with MySQL driver
- Connection via environment variable `DB_DSN`
- Automatic connection testing on startup

### HTTP Handlers
- **GET /users**: Queries all users from database
- **POST /users**: Inserts new user with JSON body
- **PUT /users/{id}**: Updates existing user by ID
- **DELETE /users/{id}**: Removes user by ID

### Error Handling
- Database connection errors
- SQL execution errors
- JSON parsing errors
- HTTP status codes for different scenarios

---

## 📦 File Overview

| File            | Purpose                              |
| --------------- | ------------------------------------ |
| `main.go`       | Go server with MySQL CRUD handlers  |
| `go.mod`        | Go module and dependencies           |
| `Dockerfile`    | Multi-stage build for the Go app     |
| `compose.yml`   | Docker Compose for MySQL + App      |
| `Makefile`      | Development and testing commands     |
| `db/init.sql`   | Database schema initialization       |

---

## 🐳 Docker Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Go App        │────│   MySQL         │
│   :8080         │    │   :3306         │
│                 │    │   (persistent)  │
└─────────────────┘    └─────────────────┘
```

- **App Container**: Runs the Go application
- **MySQL Container**: Persistent MySQL database
- **Init Script**: Automatically creates users table
- **Volume**: Data persistence between container restarts

---

## 🚨 Common Issues & Solutions

### Connection Issues
```bash
# Check if containers are running
make ps

# View application logs
make logs

# Check database connectivity
make db-cli
```

### Database Reset
```bash
# Reset everything and start fresh
make clean
make up
```

### Testing Database
```sql
-- In MySQL CLI (make db-cli)
USE learninglab;
SHOW TABLES;
SELECT * FROM users;
```

---

## 💡 Next Steps

Try these experiments:
1. **Add Validation**: Implement email format validation
2. **Error Handling**: Add custom error responses
3. **Pagination**: Add limit/offset to GET endpoint
4. **Relationships**: Create related tables (posts, comments)
5. **Transactions**: Implement multi-step operations
6. **Migrations**: Add database migration system



