package service

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/e6a5/learning/backend/04-grpc-basics/proto"
	"github.com/e6a5/learning/backend/04-grpc-basics/server/internal/models"
	"github.com/e6a5/learning/backend/04-grpc-basics/server/internal/repository"
)

// UserService implements the gRPC UserService interface
type UserService struct {
	pb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser handles unary RPC for creating a user
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Creating user: %s (%s)", req.Name, req.Email)

	user, err := s.repo.CreateUser(req.Name, req.Email)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return &pb.UserResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create user: %s", err.Error()),
		}, nil
	}

	return &pb.UserResponse{
		User:    user,
		Success: true,
		Message: "User created successfully",
	}, nil
}

// GetUser handles unary RPC for retrieving a user by ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log.Printf("Getting user with ID: %d", req.Id)

	user, err := s.repo.GetUser(req.Id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
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

// ListUsers handles unary RPC for listing users with pagination
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("Listing users: page=%d, limit=%d", req.Page, req.Limit)

	users, total, err := s.repo.ListUsers(req.Page, req.Limit)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return &pb.ListUsersResponse{
			Users: []*pb.User{},
			Total: 0,
			Page:  req.Page,
			Limit: req.Limit,
		}, fmt.Errorf("failed to list users: %w", err)
	}

	return &pb.ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// WatchUsers handles server streaming RPC for watching user creation events
func (s *UserService) WatchUsers(req *pb.WatchUsersRequest, stream pb.UserService_WatchUsersServer) error {
	log.Println("Client started watching users")

	// Create a channel for this watcher
	ch := make(chan *pb.User, 10)
	s.repo.AddWatcher(ch)

	// Remove watcher when done
	defer s.repo.RemoveWatcher(ch)

	// Send existing users first
	if err := s.sendExistingUsers(stream); err != nil {
		return fmt.Errorf("failed to send existing users: %w", err)
	}

	// Then send new users as they are created
	return s.streamNewUsers(stream, ch)
}

// BatchCreateUsers handles client streaming RPC for batch user creation
func (s *UserService) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
	log.Println("Starting batch user creation")

	requests, err := s.collectBatchRequests(stream)
	if err != nil {
		return fmt.Errorf("failed to collect batch requests: %w", err)
	}

	created, errors := s.repo.BatchCreateUsers(requests)

	log.Printf("Batch creation completed: %d created, %d errors", created, len(errors))

	return stream.SendAndClose(&pb.BatchCreateResponse{
		CreatedCount: created,
		Errors:       errors,
	})
}

// sendExistingUsers sends all existing users to the watcher stream
func (s *UserService) sendExistingUsers(stream pb.UserService_WatchUsersServer) error {
	users, _, err := s.repo.ListUsers(1, 100) // Get first 100 users
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := stream.Send(&pb.UserResponse{
			User:    user,
			Success: true,
			Message: "Existing user",
		}); err != nil {
			return err
		}
	}

	return nil
}

// streamNewUsers streams new user creation events
func (s *UserService) streamNewUsers(stream pb.UserService_WatchUsersServer, ch chan *pb.User) error {
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

// collectBatchRequests collects all requests from the batch stream
func (s *UserService) collectBatchRequests(stream pb.UserService_BatchCreateUsersServer) ([]models.CreateUserRequest, error) {
	var requests []models.CreateUserRequest

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		requests = append(requests, models.CreateUserRequest{
			Name:  req.Name,
			Email: req.Email,
		})

		log.Printf("Batch request received: %s (%s)", req.Name, req.Email)
	}

	return requests, nil
}
