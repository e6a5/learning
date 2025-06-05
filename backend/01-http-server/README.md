# 🌐 01 - HTTP Server Basics

Learn the fundamentals of **Go programming**, **package management**, and **HTTP servers** through hands-on examples. This module teaches you Go foundations while building a real HTTP API server.

---

## 🎯 What You'll Learn

### Go Fundamentals
- **Variables & Types**: strings, ints, slices, maps, structs
- **Functions**: multiple return values, error handling
- **Control Structures**: if/else, for loops, range
- **Pointers & References**: when and how to use them
- **Package System**: imports, exports, visibility

### Package Management
- **Go Modules**: initialization, versioning, dependencies
- **Third-party Packages**: finding, installing, updating
- **GitHub Integration**: importing packages directly from GitHub
- **Dependency Management**: go.mod, go.sum, semantic versioning

### HTTP Server Concepts  
- **Routing**: URL patterns, parameters, HTTP methods
- **Middleware**: request logging, authentication patterns
- **JSON APIs**: encoding/decoding, structured responses
- **Error Handling**: HTTP status codes, graceful failures

---

## 🧱 Stack

- **Go 1.23.4** - Programming language
- **gorilla/mux** - HTTP router (third-party package)
- **logrus** - Structured logging (third-party package)  
- **godotenv** - Environment variables (third-party package)

---

## 🚀 Quick Start

```bash
# Setup the learning environment
make setup

# Run the server
make run

# In another terminal, test the API
make test-api
```

Visit: http://localhost:8080

---

## 📚 Learning Path

### 1️⃣ Go Basics (Interactive Tutorial)
```bash
curl http://localhost:8080/learn/basics
```

Learn about:
- Variable declarations (`var`, `:=`, `const`)
- Data types (string, int, bool, slices, maps)
- Control structures (if/else, for, switch, range)
- Functions with multiple return values

### 2️⃣ Package Management (Live Examples)
```bash
curl http://localhost:8080/learn/packages
```

Discover:
- Standard library packages (`fmt`, `net/http`, `encoding/json`)
- Third-party packages we're using
- Popular Go packages in the ecosystem

### 3️⃣ Go Modules (Deep Dive)
```bash
curl http://localhost:8080/learn/modules
```

Understand:
- What are Go modules and why they matter
- Essential `go mod` commands
- Semantic versioning and dependency management
- Best practices for module maintenance

---

## 🔧 Package Management Hands-On

### Initialize a New Module
```bash
# Create a new Go module
go mod init github.com/yourusername/myproject

# The go.mod file defines your module
cat go.mod
```

### Add Dependencies from GitHub
```bash
# Add a specific package
go get github.com/gorilla/mux

# Add a specific version
go get github.com/sirupsen/logrus@v1.9.3

# Add the latest version
go get github.com/joho/godotenv@latest

# View what was added
go list -m all
```

### Manage Dependencies
```bash
# Clean up unused dependencies
go mod tidy

# Update all dependencies
go get -u ./...

# View dependency graph
go mod graph

# Download dependencies locally
go mod download
```

### Understanding go.mod Structure
```go
module github.com/e6a5/learning/backend/01-http-server

go 1.23.4

require (
    github.com/gorilla/mux v1.8.1     // Direct dependency
    github.com/sirupsen/logrus v1.9.3 // Direct dependency
    github.com/joho/godotenv v1.5.1   // Direct dependency
)

require (
    golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
```

---

## 🛠️ Development Commands

### Essential Go Commands
```bash
# Run without building
go run main.go

# Build a binary
go build -o server main.go

# Format code (always do this!)
go fmt ./...

# Check for common issues
go vet ./...

# Run tests
go test ./...
```

### Module Information
```bash
# Show current module info
make mod-info

# View all dependencies
make deps

# See the package management workflow
make package-workflow
```

### API Testing
```bash
# Test basic endpoints
make test-api

# Test user management
make test-users

# Test learning endpoints
make test-learning

# Full interactive demo
make demo
```

---

## 🧪 API Endpoints

### Learning Endpoints
| Endpoint | Description | What You Learn |
|----------|-------------|----------------|
| `GET /learn/basics` | Go fundamentals tutorial | Variables, types, control structures |
| `GET /learn/packages` | Package ecosystem overview | Standard library, third-party packages |
| `GET /learn/modules` | Go modules deep dive | Dependency management, versioning |

### Functional Endpoints
| Endpoint | Method | Description | Go Concepts |
|----------|--------|-------------|-------------|
| `GET /` | GET | Server info and available endpoints | JSON encoding, maps |
| `GET /health` | GET | Health check | HTTP status codes |
| `GET /users` | GET | List all users | Slices, iteration |
| `POST /users` | POST | Create a new user | JSON decoding, validation |
| `GET /users/{id}` | GET | Get user by ID | URL parameters, error handling |

---

## 💡 Go Concepts in Action

### 1. Structs and JSON Tags
```go
type User struct {
    ID       int    `json:"id"`        // JSON field mapping
    Name     string `json:"name"`      // Exported field (capitalized)
    Email    string `json:"email"`     // JSON serialization
    JoinedAt string `json:"joined_at"` // Snake case in JSON
}
```

### 2. Slices and Pointers
```go
// Slice of pointers to User structs
var users []*User

// Append to slice
users = append(users, &newUser)

// Range over slice
for _, user := range users {
    fmt.Println(user.Name)
}
```

### 3. Error Handling
```go
// Multiple return values
id, err := strconv.Atoi(idStr)
if err != nil {
    // Handle error gracefully
    sendJSONResponse(w, http.StatusBadRequest, errorResponse)
    return
}
```

### 4. Middleware Pattern
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)  // Call next handler
        duration := time.Since(start)
        // Log the request
    })
}
```

### 5. Interface{} for Generic Data
```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"` // Can hold any type
}
```

---

## 📦 Third-Party Packages Used

### gorilla/mux - HTTP Router
```go
import "github.com/gorilla/mux"

router := mux.NewRouter()
router.HandleFunc("/users/{id:[0-9]+}", getUserHandler).Methods("GET")
vars := mux.Vars(r) // Extract URL parameters
```

**Why mux over net/http?**
- URL parameter extraction
- Method-based routing
- Middleware support
- Regular expression patterns

### logrus - Structured Logging
```go
import "github.com/sirupsen/logrus"

logrus.WithFields(logrus.Fields{
    "user_id": newUser.ID,
    "method":  r.Method,
}).Info("User created")
```

**Why logrus over log?**
- Structured logging (JSON output)
- Log levels (Debug, Info, Warn, Error)
- Field-based context
- Multiple output formats

### godotenv - Environment Variables
```go
import "github.com/joho/godotenv"

godotenv.Load() // Load .env file
port := os.Getenv("PORT")
```

**Why godotenv?**
- Development environment setup
- Keep secrets out of code
- Different configs per environment
- Standard .env file format

---

## 🎭 Interactive Examples

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","email":"alice@example.com"}'
```

### Query by ID
```bash
curl http://localhost:8080/users/1
```

### Learn Go Basics
```bash
curl http://localhost:8080/learn/basics | jq '.data.variables'
```

---

## 🔍 Understanding the Code Structure

```
01-http-server/
├── main.go           # Main server implementation
├── go.mod            # Module definition and dependencies
├── go.sum            # Dependency checksums (auto-generated)
├── env.example       # Environment variable template
├── Makefile          # Development and learning commands
└── README.md         # This comprehensive guide
```

### Code Organization Patterns

**Package Declaration**
```go
package main  // Executable package
```

**Imports (grouped by type)**
```go
import (
    // Standard library
    "encoding/json"
    "net/http"
    
    // Third-party packages
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)
```

**Type Definitions**
```go
// Structs define data structures
type User struct { ... }
type Response struct { ... }
```

**Global Variables**
```go
// Package-level variables
var users []*User
var nextID int = 1
```

**Initialization**
```go
func init() {
    // Runs before main()
    // Initialize data, configuration
}
```

---

## 🎨 Best Practices Demonstrated

### 1. **Error Handling**
- Always check errors: `if err != nil`
- Return meaningful HTTP status codes
- Provide helpful error messages

### 2. **JSON API Design**
- Consistent response structure
- Proper HTTP status codes
- Content-Type headers

### 3. **Code Organization**
- One responsibility per function
- Clear function names
- Grouped imports

### 4. **Logging**
- Structured logging with context
- Different log levels
- Request/response tracking

### 5. **Environment Configuration**
- Use environment variables for config
- Provide sensible defaults
- Keep secrets out of code

---

## 🚨 Common Patterns & Gotchas

### Pointer vs Value
```go
// Good: Use pointer for modification
users = append(users, &newUser)

// Good: Use value for reading
for _, user := range users {
    fmt.Println(user.Name)
}
```

### JSON Tags Matter
```go
type User struct {
    ID   int    `json:"id"`      // lowercase in JSON
    Name string `json:"name"`    // lowercase in JSON
}
```

### Error Handling Pattern
```go
if err != nil {
    // Handle error immediately
    return
}
// Continue with success case
```

### HTTP Handler Pattern
```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse input
    // 2. Validate data
    // 3. Process business logic
    // 4. Send response
}
```

---

## 🎯 Next Steps

After mastering this module, you'll be ready for:

1. **02-mysql-crud** - Database integration
2. **03-redis-intro** - Caching and key-value storage
3. **04-grpc-basics** - High-performance RPC communication

### Exercises to Try
1. Add a `DELETE /users/{id}` endpoint
2. Add input validation for email format
3. Implement user search by name
4. Add pagination to the users list
5. Create a simple authentication middleware
6. Add unit tests for your handlers
7. Experiment with different third-party packages

### Advanced Go Topics to Explore
- **Interfaces**: Define behavior contracts
- **Channels**: Concurrent communication
- **Goroutines**: Lightweight threads
- **Context**: Request cancellation and timeouts
- **Testing**: Unit tests, table tests, mocking
- **Build Tags**: Conditional compilation

---

## 💡 Key Takeaways

✅ **Go is simple but powerful** - Few concepts, many applications  
✅ **Modules make dependency management easy** - No more GOPATH confusion  
✅ **Third-party packages extend capabilities** - Standing on shoulders of giants  
✅ **Error handling is explicit** - No hidden exceptions  
✅ **JSON APIs are straightforward** - Standard library has you covered  
✅ **Code formatting is automatic** - `go fmt` ensures consistency  

Ready to build more complex backend systems? Let's keep learning! 🚀

