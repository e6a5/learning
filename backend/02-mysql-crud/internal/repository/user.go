package repository

import (
	"database/sql"
	"fmt"

	"github.com/e6a5/learning/backend/02-mysql-crud/internal/models"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll returns all users from the database
func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

// Create creates a new user in the database
func (r *UserRepository) Create(name, email string) error {
	_, err := r.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(id string, name, email string) error {
	result, err := r.db.Exec("UPDATE users SET name=?, email=? WHERE id=?", name, email, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

// Delete deletes a user from the database
func (r *UserRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}
