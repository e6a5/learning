# 🔐 Authentication & Security

**Critical Question:** "How do you secure APIs and protect user data?"

## 🎯 Learning Objectives

### Core Concepts to Explore
- **Authentication**: Who is making the request?
- **Authorization**: What are they allowed to do?
- **Session Management**: How do you track authenticated users?
- **API Security**: How do you protect against common attacks?
- **Data Protection**: How do you secure sensitive user data?

### Security Questions We'll Answer
- "How do JWT tokens actually work?"
- "What's the difference between authentication and authorization?"
- "How do you protect against SQL injection and XSS?"
- "How do you store passwords securely?"
- "How do you implement rate limiting?"
- "What are CORS policies and why do they matter?"

---

## 🧪 Current Implementation Status: ✅ **Ready to Run**

### ✅ Implemented Features
- **JWT-based authentication** with login/logout endpoints
- **Password hashing** with bcrypt (cost factor 12)
- **Protected routes** with middleware authentication
- **Role-based authorization** (admin vs user roles)
- **Input validation** (email format, password strength)
- **SQL injection protection** with prepared statements
- **Rate limiting** (60 requests per minute per IP)
- **Security headers** (XSS, CSRF, Content-Type protection)
- **CORS configuration** for cross-origin requests
- **Pre-seeded test accounts** (admin/admin123, user/user123)

### 🗄️ Database Schema
- **users** table with roles, login tracking, account locking
- **sessions** table for tracking active JWT sessions
- **audit_logs** table for security event monitoring  
- **rate_limits** table for API rate limiting storage

### API Endpoints

| Endpoint | Method | Description | Auth Required |
|----------|--------|-------------|---------------|
| `/auth/register` | POST | Create new user account | ❌ |
| `/auth/login` | POST | Authenticate user and get JWT | ❌ |
| `/auth/logout` | POST | Invalidate user session | ✅ |
| `/auth/profile` | GET | Get current user profile | ✅ |
| `/users` | GET | List all users (admin only) | ✅ Admin |
| `/users/:id` | GET | Get specific user | ✅ Owner/Admin |
| `/users/:id` | PUT | Update user profile | ✅ Owner/Admin |
| `/users/:id` | DELETE | Delete user account | ✅ Admin |

---

## 🚀 Getting Started

### 1. Start the Database
```bash
# Start MySQL database
make up

# Check status
make ps

# View database logs
make db-logs
```

### 2. Run the Authentication Server
```bash
# Run the server locally (requires database to be running)
make run

# Or check the status endpoint
make test-status
```

### 3. Test Authentication Flow
```bash
# Test with pre-created admin account
make test-admin

# Register a new user and login
make test-auth

# Test rate limiting
make test-rate-limit
```

### 4. Database Access
```bash
# Access MySQL CLI
make db-cli

# View some sample data
SELECT * FROM users;
SELECT * FROM sessions;
```

---

## 🔍 Key Learning Points

### 1. **JWT vs Sessions**
- **JWT**: Stateless, includes user data in token
- **Sessions**: Stateful, stores data on server
- **Trade-offs**: Performance vs security vs complexity

### 2. **Password Security**
- **Never store plain text passwords**
- **Use bcrypt** for hashing (includes salt automatically)
- **Cost factor** determines computation time

### 3. **Authorization Patterns**
- **Role-based access control (RBAC)**
- **Attribute-based access control (ABAC)**
- **Resource ownership** checks

### 4. **Common Security Vulnerabilities**
- **SQL Injection**: Use prepared statements
- **XSS**: Sanitize inputs and escape outputs
- **CSRF**: Use tokens for state-changing operations
- **Rate limiting**: Prevent abuse and DoS attacks

---

## 🛡️ Security Checklist

### Authentication
- [ ] Passwords are hashed with bcrypt
- [ ] JWT tokens have appropriate expiration
- [ ] Login attempts are rate limited
- [ ] Invalid credentials don't reveal user existence

### Authorization  
- [ ] All protected routes check authentication
- [ ] Role-based access is properly enforced
- [ ] Users can only access their own data
- [ ] Admin privileges are carefully controlled

### Input Validation
- [ ] All inputs are validated and sanitized
- [ ] SQL queries use prepared statements
- [ ] File uploads are properly restricted
- [ ] Error messages don't leak sensitive info

### Infrastructure
- [ ] HTTPS is enforced in production
- [ ] Security headers are properly set
- [ ] CORS is configured correctly
- [ ] Database connections are secured

---

## 🔬 Experiments to Try

### 1. **Token Security**
- What happens if JWT secret is compromised?
- How do you rotate JWT secrets safely?
- Compare JWT vs refresh token patterns

### 2. **Attack Simulation**
- Try SQL injection attacks (safely)
- Test rate limiting effectiveness
- Simulate brute force login attempts

### 3. **Performance Impact**
- Measure authentication overhead
- Compare bcrypt cost factors
- Profile JWT vs session performance

### 4. **Architecture Variations**
- Implement OAuth2 flow
- Add multi-factor authentication
- Create API key authentication

---

## 📚 Next Steps

After mastering auth & security, explore:
- **Error Handling** (`07-error-handling/`) - Handle auth failures gracefully
- **Monitoring** (`08-monitoring/`) - Track security events and auth metrics
- **Performance** (`10-performance/`) - Optimize authentication overhead

---

## 🤝 Collaboration Notes

This implementation is designed for **learning and experimentation**. For production:
- Use established libraries (like `github.com/golang-jwt/jwt`)
- Implement proper key rotation
- Add comprehensive audit logging
- Consider OAuth2/OIDC providers
- Use HTTPS everywhere
- Implement proper session management 