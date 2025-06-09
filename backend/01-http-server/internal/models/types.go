package models

import "time"

// User represents a user in our system
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinedAt string `json:"joined_at"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate validates the create user request
func (r CreateUserRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}
	if r.Email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewUser creates a new user with generated ID and timestamp
func NewUser(name, email string, nextID int) *User {
	return &User{
		ID:       nextID,
		Name:     name,
		Email:    email,
		JoinedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}
