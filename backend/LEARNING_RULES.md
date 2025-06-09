# Backend Learning Rules - Unix Philosophy & Clean Architecture

## Core Philosophy

These rules define the foundational principles for all backend learning modules. Every new lesson must adhere to these standards to maintain consistency and quality.

## 1. Go Environment Standards (MANDATORY)

### 1.1 Go Version Requirement
- **Go Version**: 1.21 or later (REQUIRED)
- **Version Consistency**: All modules must use the same Go version
- **Version Declaration**: Specify in go.mod file explicitly

### 1.2 Module Naming Convention
- **Directory Pattern**: `XX-descriptive-name` (e.g., `01-http-server`, `05-websocket-basics`)
- **Go Module Name**: `github.com/e6a5/learning/backend/XX-descriptive-name`
- **Package Names**: Use lowercase, no underscores (e.g., `models`, `repository`, `handlers`)
- **Numbering**: Two-digit prefix for proper ordering (01, 02, ..., 10, 11, etc.)

### 1.3 Module Structure Example
```
05-websocket-basics/                    # Directory name
├── go.mod                             # module github.com/e6a5/learning/backend/05-websocket-basics
├── go.sum
├── main.go
├── compose.yml
├── README.md
├── UNIX_PHILOSOPHY.md
└── internal/
    ├── models/                        # package models
    ├── repository/                    # package repository  
    ├── handlers/                      # package handlers
    └── utils/                         # package utils
```

### 1.4 Go Module Template
```go
// go.mod file structure (MANDATORY)
module github.com/e6a5/learning/backend/XX-module-name

go 1.23.4

require (
    github.com/gorilla/mux v1.8.0
    // Add other dependencies as needed
)
```

## 2. Unix Philosophy Principles (MANDATORY)

### 2.1 "Do one thing and do it well"
- **Single Responsibility**: Each file, function, and module must have ONE clear responsibility
- **No mixing concerns**: HTTP handling, business logic, and data access must be separated
- **File size limit**: Target 50-150 lines per file (max 200 lines)
- **Function complexity**: Maximum 20-30 lines per function

### 2.2 "Write programs that work together"
- **Clean interfaces**: Define clear contracts between layers
- **Dependency injection**: NO global variables or shared state
- **Error propagation**: Consistent error handling with proper context
- **Standard patterns**: Use established patterns across all modules

### 2.3 "Write programs to handle text streams"
- **Environment configuration**: All settings via environment variables
- **Structured output**: JSON for APIs, structured logs for debugging
- **Text protocols**: Support standard protocols (HTTP, SQL, gRPC, etc.)
- **Configuration flexibility**: Deployable across different environments

## 3. Mandatory Architecture Pattern

Every backend learning module MUST follow this exact structure:

```
XX-module-name/                         # Follow naming convention
├── main.go                             # 🎯 ONLY orchestration (50-80 lines)
├── README.md                           # 📖 Module documentation  
├── go.mod                              # 📦 Go module (go 1.21, proper naming)
├── go.sum                              # 📦 Dependency checksums
├── Dockerfile                          # 🐳 Containerization (if applicable)
├── compose.yml                         # 🔧 Local development
├── internal/                  # 🔒 Private module code
│   ├── models/                # 📋 Data definitions & validation
│   │   ├── *.go              # Data structures, validation rules
│   │   └── *_test.go         # Unit tests for validation
│   ├── repository/           # 💾 Data storage operations
│   │   ├── *.go              # Storage layer (DB, cache, file)
│   │   ├── interfaces.go     # Repository contracts
│   │   └── *_test.go         # Storage layer tests
│   ├── handlers/             # 🌐 Protocol handling (HTTP/gRPC)
│   │   ├── *.go              # Request/response handling
│   │   └── *_test.go         # Handler tests
│   ├── service/              # 🔄 Business logic (if complex)
│   │   ├── *.go              # Business workflows
│   │   └── *_test.go         # Business logic tests
│   ├── middleware/           # 🛡️ Cross-cutting concerns
│   │   ├── auth.go           # Authentication
│   │   ├── logging.go        # Request logging
│   │   └── cors.go           # CORS handling
│   └── utils/                # 🔧 Common utilities
│       ├── response.go       # Response formatting
│       ├── config.go         # Configuration helpers
│       └── errors.go         # Error handling utilities
└── UNIX_PHILOSOPHY.md        # 📋 Architecture documentation
```

## 4. Code Quality Standards

### 4.1 Error Handling (MANDATORY)
```go
// ✅ CORRECT: Repository layer - detailed context
func (r *UserRepository) CreateUser(user *User) error {
    if err := r.db.Insert(user); err != nil {
        return fmt.Errorf("failed to create user %s: %w", user.Email, err)
    }
    return nil
}

// ✅ CORRECT: Handler layer - sanitized for external consumption
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    if err := h.repo.CreateUser(user); err != nil {
        log.Printf("Error creating user: %v", err)
        if isValidationError(err) {
            utils.RespondJSON(w, http.StatusBadRequest, models.APIResponse{Error: "Invalid input"})
        } else {
            utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
        }
        return
    }
}

// ❌ WRONG: Exposing internal errors
return fmt.Errorf("database connection failed: %v", dbErr)
```

### 4.2 Input Validation (MANDATORY)
```go
// ✅ CORRECT: Structured validation in models layer
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (r CreateUserRequest) Validate() error {
    if r.Name == "" {
        return &ValidationError{Field: "name", Message: "Name is required"}
    }
    if len(r.Name) < 2 {
        return &ValidationError{Field: "name", Message: "Name must be at least 2 characters"}
    }
    if !isValidEmail(r.Email) {
        return &ValidationError{Field: "email", Message: "Email format is invalid"}
    }
    return nil
}

// ❌ WRONG: Validation mixed with handlers
func CreateUser(w http.ResponseWriter, r *http.Request) {
    if user.Name == "" {  // Validation logic in handler
        http.Error(w, "Name required", 400)
        return
    }
}
```

### 4.3 Dependency Injection (MANDATORY)
```go
// ✅ CORRECT: Clean dependency injection
func main() {
    // Initialize dependencies
    db, err := database.Connect(os.Getenv("DB_URL"))
    if err != nil {
        log.Fatal(err)
    }
    
    userRepo := repository.NewUserRepository(db)
    userHandler := handlers.NewUserHandler(userRepo)
    
    // Wire up routes
    router := setupRoutes(userHandler)
    
    // Start server
    log.Fatal(http.ListenAndServe(":8080", router))
}

// ❌ WRONG: Global variables
var db *sql.DB  // Global state
var userRepo *UserRepository  // Global state

func CreateUser(w http.ResponseWriter, r *http.Request) {
    db.Insert(...)  // Direct global access
}
```

### 4.4 Import Path Standards (MANDATORY)
```go
// ✅ CORRECT: Proper import paths following module naming
package main

import (
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
    
    "github.com/e6a5/learning/backend/05-websocket-basics/internal/handlers"
    "github.com/e6a5/learning/backend/05-websocket-basics/internal/models"
    "github.com/e6a5/learning/backend/05-websocket-basics/internal/repository"
)

// ✅ CORRECT: Internal package imports
package handlers

import (
    "encoding/json"
    "net/http"
    
    "github.com/e6a5/learning/backend/05-websocket-basics/internal/models"
    "github.com/e6a5/learning/backend/05-websocket-basics/internal/repository"
    "github.com/e6a5/learning/backend/05-websocket-basics/internal/utils"
)

// ❌ WRONG: Relative imports or incorrect module paths
import (
    "./models"                    // Relative import
    "backend/models"             // Wrong module path
    "learning/models"            // Incomplete module path
)
```

## 5. Testing Requirements

### 5.1 Test Structure (MANDATORY)
Every layer must have comprehensive tests:

```
internal/
├── models/
│   ├── user.go
│   └── user_test.go      # ← Unit tests for validation
├── repository/
│   ├── user.go
│   └── user_test.go      # ← Integration tests with test DB
├── handlers/
│   ├── user.go
│   └── user_test.go      # ← HTTP tests with mock repository
└── service/
    ├── user.go
    └── user_test.go      # ← Business logic tests
```

### 5.2 Test Coverage Standards
- **Models**: 100% coverage (validation is critical)
- **Repository**: 90%+ coverage (data integrity is critical)
- **Handlers**: 85%+ coverage (API contracts are important)
- **Service**: 90%+ coverage (business logic is critical)

### 5.3 Test Categories (Use Table Testing MANDATORY)
```go
// ✅ CORRECT: Table-driven unit tests for validation
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        request models.CreateUserRequest
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid user",
            request: models.CreateUserRequest{Name: "John Doe", Email: "john@example.com"},
            wantErr: false,
        },
        {
            name:    "empty name",
            request: models.CreateUserRequest{Name: "", Email: "john@example.com"},
            wantErr: true,
            errMsg:  "Name is required",
        },
        {
            name:    "invalid email",
            request: models.CreateUserRequest{Name: "John", Email: "invalid-email"},
            wantErr: true,
            errMsg:  "Email format is invalid",
        },
        {
            name:    "short name",
            request: models.CreateUserRequest{Name: "J", Email: "j@example.com"},
            wantErr: true,
            errMsg:  "Name must be at least 2 characters",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.request.Validate()
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// ✅ CORRECT: Table-driven integration tests
func TestUserRepository_CreateUser(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    repo := repository.NewUserRepository(db)

    tests := []struct {
        name    string
        reqName string
        reqEmail string
        wantErr bool
        errMsg  string
    }{
        {
            name:     "valid user creation",
            reqName:  "John Doe",
            reqEmail: "john@example.com",
            wantErr:  false,
        },
        {
            name:     "duplicate email",
            reqName:  "Jane Doe",
            reqEmail: "john@example.com", // Same email as above
            wantErr:  true,
            errMsg:   "email already exists",
        },
        {
            name:     "invalid name",
            reqName:  "",
            reqEmail: "test@example.com",
            wantErr:  true,
            errMsg:   "Name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := repo.CreateUser(tt.reqName, tt.reqEmail)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
                assert.Equal(t, tt.reqName, user.Name)
                assert.Equal(t, tt.reqEmail, user.Email)
            }
        })
    }
}

// ✅ CORRECT: Table-driven HTTP handler tests
func TestUserHandler_CreateUser(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    string
        mockRepoSetup  func(*MockUserRepository)
        expectedStatus int
        expectedBody   string
    }{
        {
            name:        "successful creation",
            requestBody: `{"name":"John Doe","email":"john@example.com"}`,
            mockRepoSetup: func(m *MockUserRepository) {
                m.On("CreateUser", "John Doe", "john@example.com").Return(&models.User{
                    ID: 1, Name: "John Doe", Email: "john@example.com",
                }, nil)
            },
            expectedStatus: http.StatusCreated,
            expectedBody:   `"message":"User created successfully"`,
        },
        {
            name:        "invalid JSON",
            requestBody: `{"name":"John"`,
            mockRepoSetup: func(m *MockUserRepository) {
                // No mock setup needed for JSON parsing error
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `"error":"Invalid JSON"`,
        },
        {
            name:        "validation error",
            requestBody: `{"name":"","email":"john@example.com"}`,
            mockRepoSetup: func(m *MockUserRepository) {
                // No mock setup needed for validation error
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `"error":"Name is required"`,
        },
        {
            name:        "repository error",
            requestBody: `{"name":"John Doe","email":"john@example.com"}`,
            mockRepoSetup: func(m *MockUserRepository) {
                m.On("CreateUser", "John Doe", "john@example.com").Return(nil, errors.New("database error"))
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   `"error":"Internal server error"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &MockUserRepository{}
            if tt.mockRepoSetup != nil {
                tt.mockRepoSetup(mockRepo)
            }
            
            handler := handlers.NewUserHandler(mockRepo)
            req := httptest.NewRequest("POST", "/users", strings.NewReader(tt.requestBody))
            w := httptest.NewRecorder()
            
            handler.CreateUser(w, req)
            
            assert.Equal(t, tt.expectedStatus, w.Code)
            assert.Contains(t, w.Body.String(), tt.expectedBody)
            
            mockRepo.AssertExpectations(t)
        })
    }
}

// ❌ WRONG: Individual test functions (avoid this pattern)
func TestUserValidation_EmptyName(t *testing.T) {
    // Don't write separate functions for each test case
}
```

## 6. Documentation Standards

### 6.1 Module Documentation (MANDATORY)
Every module must include:

1. **README.md** with:
   - Clear purpose and learning objectives
   - Setup and running instructions
   - API documentation with examples
   - Testing instructions
   - Environment variables

2. **UNIX_PHILOSOPHY.md** with:
   - Before/after architecture comparison
   - Unix philosophy principles applied
   - Benefits achieved
   - Code quality metrics
   - Real-world impact

### 6.2 Code Documentation
```go
// ✅ CORRECT: Clear function documentation
// CreateUser creates a new user with validation and stores it in the repository.
// Returns the created user with generated ID or validation error.
func (r *UserRepository) CreateUser(name, email string) (*User, error) {
    // Implementation...
}

// ✅ CORRECT: Package documentation
// Package repository provides data access layer for user management.
// It handles database operations, connection management, and data consistency.
package repository
```

## 7. Security Standards (MANDATORY)

### 7.1 Input Security
- **SQL Injection Prevention**: Use parameterized queries ALWAYS
- **XSS Prevention**: Sanitize all user inputs
- **Input Validation**: Validate ALL inputs at the models layer
- **Rate Limiting**: Implement for all public endpoints

### 7.2 Error Security
- **No Information Leakage**: Never expose internal errors to users
- **Consistent Responses**: Same response time for found/not found
- **Audit Logging**: Log all security-relevant events

### 7.3 Configuration Security
- **Environment Variables**: Never hardcode sensitive data
- **Connection Security**: Use secure connections (TLS/SSL)
- **Access Control**: Implement proper authentication/authorization

## 8. Performance Standards

### 8.1 Database Operations
- **Connection Pooling**: Always use connection pools
- **Query Optimization**: Index frequently queried fields
- **Transaction Management**: Proper transaction boundaries
- **Resource Cleanup**: Always close resources with defer

### 8.2 Memory Management
- **Resource Leaks**: No goroutine or connection leaks
- **Buffer Management**: Proper buffer sizing
- **Garbage Collection**: Minimize GC pressure

## 9. Development Workflow

### 9.1 Module Creation Process
1. **Architecture Design**: Follow mandatory structure
2. **Interface Definition**: Define clear contracts
3. **Test-Driven Development**: Write tests first
4. **Implementation**: Follow Unix principles
5. **Documentation**: Complete all required docs
6. **Integration Testing**: End-to-end validation

### 9.2 Code Review Checklist
- [ ] **Go Environment**: Uses Go 1.21+, correct module naming
- [ ] **Unix Philosophy**: Follows all three principles
- [ ] **Architecture**: Proper layer separation, no global variables
- [ ] **Code Quality**: Comprehensive error handling, input validation
- [ ] **Testing**: Table-driven tests, meets coverage standards
- [ ] **Security**: SQL injection prevention, error sanitization
- [ ] **Documentation**: README.md and UNIX_PHILOSOPHY.md complete
- [ ] **Imports**: Correct module paths, no relative imports
- [ ] **Naming**: Follows XX-descriptive-name convention
- [ ] **Files**: compose.yml (not docker-compose.yml)

## 10. Technology Standards

### 10.1 Required Dependencies
- **HTTP Framework**: Gorilla Mux (consistency across modules)
- **Database**: Standard database/sql with drivers
- **Testing**: Standard testing package + testify for assertions
- **Environment**: Standard os package for configuration

### 10.2 Prohibited Patterns
- ❌ Global variables or shared state
- ❌ Direct database access from handlers
- ❌ Mixed concerns in single files
- ❌ Hardcoded configuration
- ❌ Exposed internal errors
- ❌ Missing input validation
- ❌ Monolithic files (>200 lines)
- ❌ Individual test functions (use table testing instead)
- ❌ Using docker-compose.yml (use compose.yml instead)
- ❌ Wrong Go version (must be 1.21+)
- ❌ Incorrect module naming (must follow XX-descriptive-name pattern)
- ❌ Relative imports or wrong import paths
- ❌ Package names with underscores or capital letters

## 11. Lesson Generation Guidelines

### 11.1 New Module Requirements
When generating a new lesson, it MUST:

1. **Follow Architecture**: Use the mandatory structure
2. **Apply Unix Principles**: Each of the three principles
3. **Include Examples**: Real-world use cases
4. **Provide Tests**: Comprehensive test coverage
5. **Document Architecture**: Before/after comparisons
6. **Show Progression**: Build on previous modules

### 11.2 Complexity Progression
- **Basic Modules**: Single protocol, simple CRUD
- **Intermediate Modules**: Multiple protocols, business logic
- **Advanced Modules**: Microservices, async processing
- **Expert Modules**: Distributed systems, performance optimization

### 11.3 Learning Objectives
Each module must clearly teach:
- **One core concept** (HTTP APIs, databases, caching, etc.)
- **Architecture patterns** (layering, separation of concerns)
- **Best practices** (error handling, validation, testing)
- **Production readiness** (security, performance, monitoring)

## Summary

These rules ensure that every backend learning module:
- **Follows Unix philosophy** consistently
- **Uses clean architecture** patterns
- **Maintains high quality** standards
- **Is production-ready** from the start
- **Is easily testable** and maintainable
- **Teaches best practices** through examples

Any deviation from these rules requires explicit justification and documentation. The goal is to create a cohesive learning experience that produces production-ready backend developers who understand both theory and practice.

**Remember**: The Unix philosophy is not just about code organization—it's about building systems that are simple, reliable, and composable. Every lesson should reinforce these principles through practical implementation. 