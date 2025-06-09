package repository

import (
	"fmt"
	"sync"

	"github.com/e6a5/learning/backend/01-http-server/internal/models"
)

// UserRepository handles user data operations
type UserRepository struct {
	users  []*models.User
	nextID int
	mutex  sync.RWMutex
}

// NewUserRepository creates a new user repository with sample data
func NewUserRepository() *UserRepository {
	repo := &UserRepository{
		users:  make([]*models.User, 0),
		nextID: 1,
	}

	// Add sample user
	sampleUser := models.NewUser("Alice Johnson", "alice@example.com", repo.nextID)
	repo.users = append(repo.users, sampleUser)
	repo.nextID++

	return repo
}

// GetAll returns all users
func (r *UserRepository) GetAll() []*models.User {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Return a copy to prevent external modification
	result := make([]*models.User, len(r.users))
	copy(result, r.users)
	return result
}

// GetByID returns a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			// Return a copy to prevent external modification
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, fmt.Errorf("user with ID %d not found", id)
}

// Create adds a new user
func (r *UserRepository) Create(name, email string) *models.User {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user := models.NewUser(name, email, r.nextID)
	r.users = append(r.users, user)
	r.nextID++

	return user
}

// Count returns the total number of users
func (r *UserRepository) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.users)
}
