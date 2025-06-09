# Unix Philosophy Refactoring - 03-Redis-Intro

## The Problem: Monolithic main.go (204 lines)

**Before**: Single `main.go` with 204 lines violating:
- **Single Responsibility Principle** - doing everything in one file
- **Unix Philosophy** - "do one thing and do it well"
- **Global Variables** - shared Redis client and context
- **Separation of Concerns** - HTTP, Redis, and business logic mixed together

## The Solution: Modular Architecture

**After**: Broken down into focused, single-purpose modules:

```
03-redis-intro/
├── main.go (54 lines)              # 🎯 Application orchestration only
├── internal/
│   ├── models/                     # 📋 Data definitions & validation
│   │   └── cache.go (67 lines)     
│   ├── repository/                 # 💾 Redis operations
│   │   └── cache.go (102 lines)    
│   ├── handlers/                   # 🌐 HTTP request handling
│   │   └── cache.go (139 lines)    
│   └── utils/                      # 🔧 Common utilities
│       └── response.go (25 lines)  
```

## Unix Philosophy Applied

### 1. **"Do one thing and do it well"**

| Module | Single Responsibility |
|--------|----------------------|
| `main.go` | **Application orchestration** - Redis setup, dependency wiring, server startup |
| `models/cache.go` | **Data definitions** - Redis models, validation rules, business constraints |
| `repository/cache.go` | **Redis operations** - Cache commands, connection management, error handling |
| `handlers/cache.go` | **HTTP endpoints** - Request parsing, response formatting, HTTP status codes |
| `utils/response.go` | **Response utilities** - JSON responses, environment variables |

### 2. **"Write programs that work together"**

- **Clean interfaces** between layers (handler → repository → Redis)
- **Dependency injection** eliminating global variables
- **Error propagation** with proper context and wrapping
- **Standard HTTP patterns** following REST conventions

### 3. **"Write programs to handle text streams"**

- **JSON API** with standard request/response formats
- **Environment-based configuration** (REDIS_ADDR, REDIS_PASSWORD)
- **Structured error messages** with field-level validation
- **Cache key patterns** supporting wildcard searches

## Benefits Achieved

### 📏 **File Size Reduction**
- **Before**: 1 file × 204 lines = **monolithic structure**
- **After**: 5 files × ~77 lines avg = **focused & manageable**

### 🛡️ **Architectural Improvements**
- **Before**: Global Redis client accessible everywhere
- **After**: Dependency injection with clean interfaces
- **Connection Management**: Centralized in repository layer
- **Error Handling**: Comprehensive with proper context

### 🧪 **Testability**
- **Before**: Testing required Redis connection for everything
- **After**: Each layer can be tested independently
- **Repository Testing**: Mock Redis interfaces
- **Handler Testing**: Mock repository implementations
- **Validation Testing**: Unit test validation rules

### 🔧 **Maintainability**
- **Before**: Change cache logic? Hunt through HTTP handlers
- **After**: Change cache logic? Modify `repository/cache.go` (102 lines)
- **Before**: Change validation? Modify mixed HTTP/Redis code
- **After**: Change validation? Modify `models/cache.go` (67 lines)

### 🏗️ **Extensibility**
- **Before**: Adding Redis clustering = rewriting global client code
- **After**: Adding Redis clustering = modify repository initialization
- **Before**: Adding middleware = modifying HTTP handlers
- **After**: Adding middleware = add to router setup

## Code Quality Metrics

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| **Largest file** | 204 lines | 139 lines | 32% reduction |
| **Average file size** | 204 lines | 77 lines | 62% reduction |
| **Global state** | 2 globals | 0 globals | ✅ |
| **Error handling** | Basic | Comprehensive | Significant |
| **Input validation** | Manual | Structured | ✅ |
| **Testability** | Difficult | Easy | ✅ |

## Redis Integration Improvements

### **Before (Global State)**
```go
var rdb *redis.Client
var ctx = context.Background()

func getValue(w http.ResponseWriter, r *http.Request) {
    val, err := rdb.Get(ctx, key).Result()  // Global variables
    if err == redis.Nil {
        respondJSON(w, http.StatusNotFound, Response{Error: "Key not found"})
        return
    }
    // Mixed HTTP and Redis logic
}
```

### **After (Clean Architecture)**
```go
// Repository layer - single responsibility for Redis
func (r *CacheRepository) Get(key string) (*models.KeyValue, error) {
    val, err := r.client.Get(r.ctx, key).Result()
    if err == redis.Nil {
        return nil, fmt.Errorf("key not found: %s", key)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get key %s: %w", key, err)
    }
    return models.NewKeyValue(key, val, 0), nil
}

// Handler layer - single responsibility for HTTP
func (h *CacheHandler) GetValue(w http.ResponseWriter, r *http.Request) {
    key := mux.Vars(r)["key"]
    
    kv, err := h.repo.Get(key)  // Clean interface
    if err != nil {
        // Proper error handling with logging
        log.Printf("Error getting key %s: %v", key, err)
        if strings.Contains(err.Error(), "key not found") {
            utils.RespondJSON(w, http.StatusNotFound, models.APIResponse{Error: "Key not found"})
        } else {
            utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
        }
        return
    }
    
    utils.RespondJSON(w, http.StatusOK, models.APIResponse{Data: kv})
}
```

## Architecture Philosophy

### **Before (Global State Monolith)**
```
main.go
├── Global Redis client
├── Global context
├── HTTP handlers mixed with Redis calls
├── Manual JSON handling
├── Basic error handling
└── All concerns in one place
```
☠️ **Everything tightly coupled with global state**

### **After (Clean Architecture)**
```
main.go (orchestrator)
├── initializes Redis connection
├── wires dependencies
└── starts server

internal/ (focused layers)
├── models/ (validation + data structures)
├── repository/ (Redis + cache operations + error handling)
├── handlers/ (HTTP + request/response + status codes)
└── utils/ (common + response formatting)
```
✅ **Each layer has single responsibility, no global state**

## Real-World Impact

### 🚀 **Developer Experience**
```bash
# Find Redis operations
vim internal/repository/cache.go    # 102 lines, Redis commands

# Find validation logic  
vim internal/models/cache.go        # 67 lines, validation rules

# Find HTTP endpoint logic
vim internal/handlers/cache.go      # 139 lines, HTTP handling
```

### 🧪 **Testing Strategy**
```bash
# Test Redis operations with mocks
go test ./internal/repository

# Test validation independently
go test ./internal/models

# Test HTTP handlers with mock repository
go test ./internal/handlers
```

### 🛡️ **Production Benefits**
- **Connection pooling** managed in repository layer
- **Graceful Redis reconnection** handled centrally
- **Environment-based configuration** for different deployments
- **Proper error sanitization** preventing information leakage

## Next Steps for Production

With this modular structure, we can now add:

```
03-redis-intro/
├── internal/
│   ├── models/
│   │   ├── cache.go
│   │   └── cache_test.go       # ← Validation tests
│   ├── repository/
│   │   ├── cache.go
│   │   └── cache_test.go       # ← Redis integration tests
│   ├── handlers/
│   │   ├── cache.go
│   │   └── cache_test.go       # ← HTTP handler tests
│   └── middleware/             # ← Add rate limiting, auth
│       └── ratelimit.go
```

## Summary

**From 1 × 204-line monolith → 5 × ~77-line focused modules**

This refactoring demonstrates **Unix philosophy in caching applications**:

- **Does one thing well** ✅ - Each module has clear purpose
- **Works with others** ✅ - Clean interfaces between layers  
- **Handles text streams** ✅ - JSON API with environment config
- **Eliminates global state** ✅ - Dependency injection throughout
- **Is testable in isolation** ✅ - Each layer independently testable

The result is **scalable, maintainable, production-ready cache API** following time-tested principles. 