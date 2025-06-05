package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/e6a5/learning/backend/04-grpc-basics/proto"
	"google.golang.org/grpc"
)

// In-memory user storage
type UserStore struct {
	mu       sync.RWMutex
	users    map[int32]*pb.User
	nextID   int32
	watchers []chan *pb.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:    make(map[int32]*pb.User),
		nextID:   1,
		watchers: make([]chan *pb.User, 0),
	}
}

func (s *UserStore) CreateUser(name, email string) *pb.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &pb.User{
		Id:        s.nextID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().Unix(),
	}

	s.users[s.nextID] = user
	s.nextID++

	// Notify watchers
	for _, watcher := range s.watchers {
		select {
		case watcher <- user:
		default:
			// Channel is full, skip
		}
	}

	return user
}

func (s *UserStore) GetUser(id int32) (*pb.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[id]
	return user, exists
}

func (s *UserStore) ListUsers(page, limit int32) ([]*pb.User, int32) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*pb.User
	for _, user := range s.users {
		users = append(users, user)
	}

	// Simple pagination
	start := (page - 1) * limit
	end := start + limit
	total := int32(len(users))

	if start >= total {
		return []*pb.User{}, total
	}
	if end > total {
		end = total
	}

	return users[start:end], total
}

func (s *UserStore) AddWatcher(ch chan *pb.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watchers = append(s.watchers, ch)
}

func (s *UserStore) RemoveWatcher(ch chan *pb.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, watcher := range s.watchers {
		if watcher == ch {
			s.watchers = append(s.watchers[:i], s.watchers[i+1:]...)
			close(ch)
			break
		}
	}
}

// UserService implements the gRPC service
type UserService struct {
	pb.UnimplementedUserServiceServer
	store *UserStore
}

func NewUserService() *UserService {
	return &UserService{
		store: NewUserStore(),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Creating user: %s (%s)", req.Name, req.Email)

	if req.Name == "" || req.Email == "" {
		return &pb.UserResponse{
			Success: false,
			Message: "Name and email are required",
		}, nil
	}

	user := s.store.CreateUser(req.Name, req.Email)

	return &pb.UserResponse{
		User:    user,
		Success: true,
		Message: "User created successfully",
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log.Printf("Getting user with ID: %d", req.Id)

	user, exists := s.store.GetUser(req.Id)
	if !exists {
		return &pb.UserResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	return &pb.UserResponse{
		User:    user,
		Success: true,
		Message: "User found",
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("Listing users: page=%d, limit=%d", req.Page, req.Limit)

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total := s.store.ListUsers(req.Page, req.Limit)

	return &pb.ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (s *UserService) WatchUsers(req *pb.WatchUsersRequest, stream pb.UserService_WatchUsersServer) error {
	log.Println("Client started watching users")

	// Create a channel for this watcher
	ch := make(chan *pb.User, 10)
	s.store.AddWatcher(ch)

	// Remove watcher when done
	defer s.store.RemoveWatcher(ch)

	// Send existing users first
	users, _ := s.store.ListUsers(1, 100) // Get first 100 users
	for _, user := range users {
		if err := stream.Send(&pb.UserResponse{
			User:    user,
			Success: true,
			Message: "Existing user",
		}); err != nil {
			return err
		}
	}

	// Then send new users as they are created
	for {
		select {
		case user := <-ch:
			if err := stream.Send(&pb.UserResponse{
				User:    user,
				Success: true,
				Message: "New user created",
			}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			log.Println("Client stopped watching users")
			return nil
		}
	}
}

func (s *UserService) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
	log.Println("Starting batch user creation")

	var created int32
	var errors []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Client finished sending
			return stream.SendAndClose(&pb.BatchCreateResponse{
				CreatedCount: created,
				Errors:       errors,
			})
		}
		if err != nil {
			return err
		}

		if req.Name == "" || req.Email == "" {
			errors = append(errors, fmt.Sprintf("Invalid user: name='%s', email='%s'", req.Name, req.Email))
			continue
		}

		s.store.CreateUser(req.Name, req.Email)
		created++
		log.Printf("Batch created user: %s (%s)", req.Name, req.Email)
	}
}

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userService := NewUserService()

	pb.RegisterUserServiceServer(grpcServer, userService)

	log.Printf("🚀 gRPC Server running on port %d", port)
	log.Println("Available services:")
	log.Println("  - CreateUser (unary)")
	log.Println("  - GetUser (unary)")
	log.Println("  - ListUsers (unary)")
	log.Println("  - WatchUsers (server streaming)")
	log.Println("  - BatchCreateUsers (client streaming)")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
