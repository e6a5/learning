package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"

	_ "github.com/go-sql-driver/mysql"
)

// 🔐 Configuration
const (
	JWTSecret   = "your-secret-key-change-in-production"
	BCryptCost  = 12
	TokenExpiry = 24 * time.Hour
	ServerPort  = ":8080"
)

func getDatabaseDSN() string {
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		return dsn
	}
	return "user:pass@tcp(localhost:3306)/authlab?parseTime=true"
}

// 📊 Data Structures
type User struct {
	ID                  int        `json:"id"`
	Username            string     `json:"username"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"` // Never send to client
	Role                string     `json:"role"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	IsActive            bool       `json:"is_active"`
	LastLogin           *time.Time `json:"last_login,omitempty"`
	FailedLoginAttempts int        `json:"-"`
	LockedUntil         *time.Time `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 🏗️ Application Structure
type AuthServer struct {
	db      *sql.DB
	limiter map[string]*rate.Limiter
}

// 🔧 Helper Functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BCryptCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(user User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func validatePassword(password string) bool {
	return len(password) >= 8
}

// 🛡️ Security Middleware
func (s *AuthServer) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *AuthServer) rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if s.limiter[ip] == nil {
			s.limiter[ip] = rate.NewLimiter(rate.Every(time.Minute), 60) // 60 requests per minute
		}

		if !s.limiter[ip].Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *AuthServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AuthServer) adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 📝 Database Operations
func (s *AuthServer) createUser(user RegisterRequest) (*User, error) {
	// Validate input
	if !validateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format")
	}
	if !validatePassword(user.Password) {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Insert user
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES (?, ?, ?)
	`
	result, err := s.db.Exec(query, user.Username, user.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Return created user
	return s.getUserByID(int(id))
}

func (s *AuthServer) getUserByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at, 
		       is_active, last_login, failed_login_attempts, locked_until
		FROM users WHERE username = ?
	`
	var user User
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive,
		&user.LastLogin, &user.FailedLoginAttempts, &user.LockedUntil,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthServer) getUserByID(id int) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at, 
		       is_active, last_login, failed_login_attempts, locked_until
		FROM users WHERE id = ?
	`
	var user User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive,
		&user.LastLogin, &user.FailedLoginAttempts, &user.LockedUntil,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthServer) updateLastLogin(userID int) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = ?`
	_, err := s.db.Exec(query, userID)
	return err
}

// 🔐 HTTP Handlers
func (s *AuthServer) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := s.createUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "Username or email already exists", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func (s *AuthServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := s.getUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !checkPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		http.Error(w, "Account is disabled", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(*user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Update last login
	s.updateLastLogin(user.ID)

	response := LoginResponse{
		Token:   token,
		User:    *user,
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *AuthServer) profileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}

	user, err := s.getUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *AuthServer) usersHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, username, email, role, created_at, updated_at, is_active, last_login
		FROM users ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Role,
			&user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.LastLogin,
		)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

func (s *AuthServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"message":   "🔐 Authentication & Security Server is running",
		"endpoints": map[string]string{
			"POST /auth/register": "Create new user account",
			"POST /auth/login":    "Authenticate user and get JWT",
			"GET /auth/profile":   "Get current user profile (auth required)",
			"GET /users":          "List all users (admin only)",
		},
	})
}

// 🚀 Server Setup
func initDB() (*sql.DB, error) {
	dsn := getDatabaseDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("✅ Connected to MySQL database at %s", dsn)
	return db, nil
}

func main() {
	log.Println("🔐 Starting Authentication & Security Server...")

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// Create server
	server := &AuthServer{
		db:      db,
		limiter: make(map[string]*rate.Limiter),
	}

	// Setup routes
	r := mux.NewRouter()

	// Apply security middleware to all routes
	r.Use(server.securityHeaders)
	r.Use(server.rateLimiter)

	// Public routes
	r.HandleFunc("/", server.statusHandler).Methods("GET")
	r.HandleFunc("/auth/register", server.registerHandler).Methods("POST")
	r.HandleFunc("/auth/login", server.loginHandler).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/auth").Subrouter()
	protected.Use(server.authMiddleware)
	protected.HandleFunc("/profile", server.profileHandler).Methods("GET")

	// Admin routes
	admin := r.PathPrefix("/users").Subrouter()
	admin.Use(server.authMiddleware)
	admin.Use(server.adminOnly)
	admin.HandleFunc("", server.usersHandler).Methods("GET")

	log.Printf("🚀 Server starting on port %s", ServerPort)
	log.Println("📚 Available endpoints:")
	log.Println("  GET  /                - Server status")
	log.Println("  POST /auth/register   - Create user account")
	log.Println("  POST /auth/login      - Authenticate user")
	log.Println("  GET  /auth/profile    - Get user profile (auth required)")
	log.Println("  GET  /users           - List users (admin only)")

	if err := http.ListenAndServe(ServerPort, r); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}
