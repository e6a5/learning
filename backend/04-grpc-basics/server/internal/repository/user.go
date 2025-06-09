package repository

import (
	"fmt"
	"sync"

	pb "github.com/e6a5/learning/backend/04-grpc-basics/proto"
	"github.com/e6a5/learning/backend/04-grpc-basics/server/internal/models"
)

// UserRepository handles user storage operations
type UserRepository struct {
	mu       sync.RWMutex
	users    map[int32]*pb.User
	nextID   int32
	watchers []chan *pb.User
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:    make(map[int32]*pb.User),
		nextID:   1,
		watchers: make([]chan *pb.User, 0),
	}
}

// CreateUser creates a new user with validation
func (r *UserRepository) CreateUser(name, email string) (*pb.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, err := models.NewUser(r.nextID, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.users[r.nextID] = user
	r.nextID++

	// Notify watchers
	r.notifyWatchers(user)

	return user, nil
}

// GetUser retrieves a user by ID
func (r *UserRepository) GetUser(id int32) (*pb.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found: %d", id)
	}

	return user, nil
}

// ListUsers returns paginated users
func (r *UserRepository) ListUsers(page, limit int32) ([]*pb.User, int32, error) {
	normalizedPage, normalizedLimit, err := models.NormalizeListRequest(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid list request: %w", err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var users []*pb.User
	for _, user := range r.users {
		users = append(users, user)
	}

	// Simple pagination
	start := (normalizedPage - 1) * normalizedLimit
	end := start + normalizedLimit
	total := int32(len(users))

	if start >= total {
		return []*pb.User{}, total, nil
	}
	if end > total {
		end = total
	}

	return users[start:end], total, nil
}

// AddWatcher adds a new watcher for user creation events
func (r *UserRepository) AddWatcher(ch chan *pb.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.watchers = append(r.watchers, ch)
}

// RemoveWatcher removes a watcher
func (r *UserRepository) RemoveWatcher(ch chan *pb.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, watcher := range r.watchers {
		if watcher == ch {
			r.watchers = append(r.watchers[:i], r.watchers[i+1:]...)
			close(ch)
			break
		}
	}
}

// GetUserCount returns the total number of users
func (r *UserRepository) GetUserCount() int32 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int32(len(r.users))
}

// BatchCreateUsers creates multiple users and returns results
func (r *UserRepository) BatchCreateUsers(requests []models.CreateUserRequest) (int32, []string) {
	var created int32
	var errors []string

	for _, req := range requests {
		if err := req.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("Invalid user: name='%s', email='%s' - %s", req.Name, req.Email, err.Error()))
			continue
		}

		_, err := r.CreateUser(req.Name, req.Email)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to create user: name='%s', email='%s' - %s", req.Name, req.Email, err.Error()))
			continue
		}

		created++
	}

	return created, errors
}

// notifyWatchers sends user creation events to all watchers
func (r *UserRepository) notifyWatchers(user *pb.User) {
	for _, watcher := range r.watchers {
		select {
		case watcher <- user:
		default:
			// Channel is full, skip to avoid blocking
		}
	}
}
