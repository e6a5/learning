package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateUserRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateUserRequest
		expectError bool
		errorField  string
		errorMsg    string
	}{
		{
			name:        "valid user request",
			request:     CreateUserRequest{Name: "John Doe", Email: "john@example.com"},
			expectError: false,
		},
		{
			name:        "empty name",
			request:     CreateUserRequest{Name: "", Email: "john@example.com"},
			expectError: true,
			errorField:  "name",
			errorMsg:    "name is required",
		},
		{
			name:        "whitespace only name",
			request:     CreateUserRequest{Name: "   ", Email: "john@example.com"},
			expectError: true,
			errorField:  "name",
			errorMsg:    "name is required",
		},
		{
			name:        "name too long",
			request:     CreateUserRequest{Name: strings.Repeat("a", 101), Email: "john@example.com"},
			expectError: true,
			errorField:  "name",
			errorMsg:    "name must be 100 characters or less",
		},
		{
			name:        "empty email",
			request:     CreateUserRequest{Name: "John Doe", Email: ""},
			expectError: true,
			errorField:  "email",
			errorMsg:    "email is required",
		},
		{
			name:        "invalid email format",
			request:     CreateUserRequest{Name: "John Doe", Email: "invalid-email"},
			expectError: true,
			errorField:  "email",
			errorMsg:    "email format is invalid",
		},
		{
			name:        "invalid email - missing @",
			request:     CreateUserRequest{Name: "John Doe", Email: "johndoe.com"},
			expectError: true,
			errorField:  "email",
			errorMsg:    "email format is invalid",
		},
		{
			name:        "invalid email - missing domain",
			request:     CreateUserRequest{Name: "John Doe", Email: "john@"},
			expectError: true,
			errorField:  "email",
			errorMsg:    "email format is invalid",
		},
		{
			name:        "valid email with numbers",
			request:     CreateUserRequest{Name: "John Doe", Email: "john123@example.com"},
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			request:     CreateUserRequest{Name: "John Doe", Email: "john@mail.example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateUserRequest(tt.request)

			if tt.expectError {
				require.Error(t, err, "Expected validation error")

				// Check if it's our custom validation error
				validationErr, ok := err.(UserValidationError)
				require.True(t, ok, "Expected UserValidationError type")

				assert.Equal(t, tt.errorField, validationErr.Field)
				assert.Contains(t, validationErr.Message, tt.errorMsg)
			} else {
				assert.NoError(t, err, "Expected no validation error")
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		request  CreateUserRequest
		id       int
		validate func(t *testing.T, user User)
	}{
		{
			name:    "create user with basic info",
			request: CreateUserRequest{Name: "John Doe", Email: "john@example.com"},
			id:      1,
			validate: func(t *testing.T, user User) {
				assert.Equal(t, 1, user.ID)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john@example.com", user.Email)
				assert.NotEmpty(t, user.JoinedAt)
			},
		},
		{
			name:    "trims whitespace from name",
			request: CreateUserRequest{Name: "  John Doe  ", Email: "john@example.com"},
			id:      2,
			validate: func(t *testing.T, user User) {
				assert.Equal(t, "John Doe", user.Name) // Should be trimmed
			},
		},
		{
			name:    "normalizes email to lowercase",
			request: CreateUserRequest{Name: "John Doe", Email: "JOHN@EXAMPLE.COM"},
			id:      3,
			validate: func(t *testing.T, user User) {
				assert.Equal(t, "john@example.com", user.Email) // Should be lowercase
			},
		},
		{
			name:    "trims whitespace from email",
			request: CreateUserRequest{Name: "John Doe", Email: "  john@example.com  "},
			id:      4,
			validate: func(t *testing.T, user User) {
				assert.Equal(t, "john@example.com", user.Email) // Should be trimmed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.request, tt.id)
			tt.validate(t, user)
		})
	}
}

func TestUser_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected bool
	}{
		{
			name:     "zero value user is empty",
			user:     User{},
			expected: true,
		},
		{
			name:     "user with only ID is not empty",
			user:     User{ID: 1},
			expected: false,
		},
		{
			name:     "user with only name is not empty",
			user:     User{Name: "John"},
			expected: false,
		},
		{
			name:     "user with only email is not empty",
			user:     User{Email: "john@example.com"},
			expected: false,
		},
		{
			name:     "fully populated user is not empty",
			user:     User{ID: 1, Name: "John", Email: "john@example.com", JoinedAt: "2023-01-01"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_SanitizeForDisplay(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		JoinedAt: "2023-01-01 10:00:00",
	}

	sanitized := user.SanitizeForDisplay()

	// Currently, we return the user as-is since no sensitive data
	assert.Equal(t, user, sanitized)
}

func TestUserValidationError_Error(t *testing.T) {
	err := UserValidationError{
		Field:   "email",
		Message: "is required",
	}

	expected := "email: is required"
	assert.Equal(t, expected, err.Error())
}

// Benchmark tests for performance measurement
func BenchmarkValidateCreateUserRequest(b *testing.B) {
	request := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCreateUserRequest(request)
	}
}

func BenchmarkNewUser(b *testing.B) {
	request := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewUser(request, i)
	}
}

// Test helper functions
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"john@example.com", true},
		{"john.doe@example.com", true},
		{"john+tag@example.com", true},
		{"john@mail.example.com", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"john@", false},
		{"john.example.com", false},
		{"", false},
		{"john@.com", false},
		{"john@example.", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.expected, result, "Email: %s", tt.email)
		})
	}
}
