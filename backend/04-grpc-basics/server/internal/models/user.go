package models

import (
	"fmt"
	"strings"
	"time"

	pb "github.com/e6a5/learning/backend/04-grpc-basics/proto"
)

// CreateUserRequest represents the validated user creation request
type CreateUserRequest struct {
	Name  string
	Email string
}

// ListUsersRequest represents the validated user list request
type ListUsersRequest struct {
	Page  int32
	Limit int32
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate validates user creation input
func (r CreateUserRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}
	if len(r.Name) < 2 {
		return &ValidationError{Field: "name", Message: "Name must be at least 2 characters"}
	}
	if len(r.Name) > 100 {
		return &ValidationError{Field: "name", Message: "Name must be less than 100 characters"}
	}

	if r.Email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	if !isValidEmail(r.Email) {
		return &ValidationError{Field: "email", Message: "Email format is invalid"}
	}

	return nil
}

// Validate validates list users request
func (r ListUsersRequest) Validate() error {
	if r.Page < 0 {
		return &ValidationError{Field: "page", Message: "Page must be non-negative"}
	}
	if r.Limit < 0 {
		return &ValidationError{Field: "limit", Message: "Limit must be non-negative"}
	}
	if r.Limit > 1000 {
		return &ValidationError{Field: "limit", Message: "Limit must be less than 1000"}
	}

	return nil
}

// NewUser creates a new protobuf User with validation
func NewUser(id int32, name, email string) (*pb.User, error) {
	req := CreateUserRequest{Name: name, Email: email}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        id,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().Unix(),
	}, nil
}

// NormalizeListRequest normalizes and validates list request
func NormalizeListRequest(page, limit int32) (int32, int32, error) {
	req := ListUsersRequest{Page: page, Limit: limit}

	// Apply defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if err := req.Validate(); err != nil {
		return 0, 0, err
	}

	return req.Page, req.Limit, nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if len(email) < 5 || len(email) > 254 {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]
	if len(local) == 0 || len(local) > 64 {
		return false
	}
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}
