package models

import (
	"regexp"
	"strings"
	"time"
)

// User represents a user in our system
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinedAt string `json:"joined_at"`
}

// CreateUserRequest represents the payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserValidationError represents validation errors
type UserValidationError struct {
	Field   string
	Message string
}

func (e UserValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// ValidateCreateUserRequest validates a user creation request
func ValidateCreateUserRequest(req CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return UserValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	if len(req.Name) > 100 {
		return UserValidationError{
			Field:   "name",
			Message: "name must be 100 characters or less",
		}
	}

	if strings.TrimSpace(req.Email) == "" {
		return UserValidationError{
			Field:   "email",
			Message: "email is required",
		}
	}

	if !isValidEmail(req.Email) {
		return UserValidationError{
			Field:   "email",
			Message: "email format is invalid",
		}
	}

	return nil
}

// NewUser creates a new user with generated ID and timestamp
func NewUser(req CreateUserRequest, id int) User {
	return User{
		ID:       id,
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.TrimSpace(strings.ToLower(req.Email)),
		JoinedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsEmpty checks if a user is considered empty/zero-value
func (u User) IsEmpty() bool {
	return u.ID == 0 && u.Name == "" && u.Email == ""
}

// SanitizeForDisplay removes sensitive information for public display
func (u User) SanitizeForDisplay() User {
	// In a real system, you might remove sensitive fields
	// For now, we just return the user as-is since no sensitive data
	return u
}
