# 🧪 05 - Testing Basics

**Concept-Driven Testing** — How do you ensure validation logic works reliably?

This module explores testing **user validation functions** to answer the fundamental question: *"How do you know your input validation works correctly before it reaches your database or APIs?"*

---

## 🎯 Central Question

> **"How do you ensure validation logic works reliably?"**

### What We're Testing
- User input validation (name, email)
- Edge cases (empty strings, invalid formats)  
- Performance of validation functions
- Error handling and custom error types

---

## 🧠 What We Built

**Problem**: We have validation functions but don't know if they handle all edge cases correctly.

**Solution**: Comprehensive unit tests with 100% coverage that verify every validation scenario.

**Result**: Confidence that validation logic catches all invalid inputs before they reach our system.

---

## 🔬 What We Actually Test

### ✅ **Validation Functions**
- `ValidateCreateUserRequest()` - Input validation
- `NewUser()` - User creation with normalization
- `isValidEmail()` - Email format validation
- Error handling with custom error types

### ✅ **Test Scenarios Covered**
- **Valid inputs**: Proper name and email formats
- **Invalid names**: Empty, whitespace-only, too long (>100 chars)
- **Invalid emails**: Missing @, invalid format, empty
- **Edge cases**: Whitespace trimming, email normalization
- **Performance**: Benchmark validation speed

---

## 🧱 Stack

- **Go testing** - Built-in testing framework (`go test`)
- **testify** - Better assertions (`assert.Equal`, `require.Error`)
- **Table-driven tests** - Test multiple scenarios efficiently
- **Benchmarks** - Measure validation performance

---

## 🚀 Quick Start

```bash
# Run all tests
make test

# Run tests with coverage report
make coverage

# Run performance benchmarks
make bench

# See all available commands
make help
```

---

## 🧪 Test Results

### Test Coverage: 100%
```bash
$ make coverage
✅ TestValidateCreateUserRequest (10 scenarios)
✅ TestNewUser (4 scenarios)  
✅ TestUser_IsEmpty (5 scenarios)
✅ TestIsValidEmail (11 scenarios)

Coverage: 100.0% of statements
```

### Performance Benchmarks
```bash
$ make bench
BenchmarkValidateCreateUserRequest-12    213238    5642 ns/op
BenchmarkNewUser-12                     4304302     266 ns/op
```

---

## 🔧 Test Examples

### Table-Driven Validation Tests
```go
func TestValidateCreateUserRequest(t *testing.T) {
    tests := []struct {
        name        string
        request     CreateUserRequest
        expectError bool
        errorField  string
    }{
        {
            name:        "valid user request",
            request:     CreateUserRequest{Name: "John", Email: "john@test.com"},
            expectError: false,
        },
        {
            name:        "empty name",
            request:     CreateUserRequest{Name: "", Email: "john@test.com"},
            expectError: true,
            errorField:  "name",
        },
        // 8 more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateCreateUserRequest(tt.request)
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Email Validation Tests
```go
func TestIsValidEmail(t *testing.T) {
    tests := []struct {
        email    string
        expected bool
    }{
        {"john@example.com", true},
        {"invalid-email", false},
        {"@example.com", false},
        // 8 more cases...
    }
    
    for _, tt := range tests {
        result := isValidEmail(tt.email)
        assert.Equal(t, tt.expected, result)
    }
}
```

---

## 🎯 What We Learned

### ✅ **Testing Benefits Proven**
- **Caught edge cases** we hadn't considered (whitespace-only names)
- **Documented behavior** through test examples
- **Enabled refactoring** with confidence
- **Measured performance** (validation takes ~5μs per request)

### ✅ **Testing Patterns Used**
- **Table-driven tests** for multiple scenarios
- **testify assertions** for clear failure messages
- **Benchmark tests** for performance measurement
- **Custom error types** for structured error handling

---

## 📁 Project Structure

```
05-testing-basics/
├── models/
│   ├── user.go          # Validation functions we're testing
│   └── user_test.go     # 30+ test cases with 100% coverage
├── Makefile            # Test automation commands
├── go.mod              # Dependencies (testify only)
└── coverage.out        # Generated coverage report
```

---

## 🤖 AI Collaboration Opportunities

### Extend Testing
**Ask AI**: *"What validation edge cases am I missing? Should I test unicode characters in names?"*

### Improve Performance  
**Ask AI**: *"This email regex takes 5μs - is that slow? How can I optimize it?"*

### Add More Validation
**Ask AI**: *"Help me add password validation with complexity requirements and tests"*

---

## 💡 Next Questions This Raises

After proving validation logic works reliably:
- **"How do I test HTTP endpoints that use this validation?"** → HTTP endpoint testing
- **"How do I test database operations with this validated data?"** → Integration testing
- **"How do I test the full user creation flow?"** → End-to-end testing

---

## 🎯 Success Metrics

✅ **100% test coverage** on validation functions  
✅ **30+ test scenarios** covering edge cases  
✅ **Performance benchmarks** showing validation speed  
✅ **Zero validation bugs** in manual testing  
✅ **Clear error messages** for invalid inputs

**Ready to ensure your validation logic works perfectly?** Run `make test` and see all scenarios pass! 🧪 