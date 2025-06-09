package models

// User represents a user in the database
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
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

// Validate validates the update user request
func (r UpdateUserRequest) Validate() error {
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
