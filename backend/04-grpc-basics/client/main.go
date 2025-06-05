package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/e6a5/learning/backend/04-grpc-basics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx := context.Background()

	log.Println("🔗 Connected to gRPC server")
	log.Println("🧪 Running gRPC client examples...")

	// 1. Test Unary RPC - Create User
	log.Println("\n1️⃣ Testing Unary RPC - CreateUser")
	createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	log.Printf("✅ Created user: %s (ID: %d)", createResp.User.Name, createResp.User.Id)

	// Create another user
	createResp2, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Bob Smith",
		Email: "bob@example.com",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	log.Printf("✅ Created user: %s (ID: %d)", createResp2.User.Name, createResp2.User.Id)

	// 2. Test Unary RPC - Get User
	log.Println("\n2️⃣ Testing Unary RPC - GetUser")
	getResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	if getResp.Success {
		log.Printf("✅ Found user: %s (%s)", getResp.User.Name, getResp.User.Email)
	}

	// 3. Test Unary RPC - List Users
	log.Println("\n3️⃣ Testing Unary RPC - ListUsers")
	listResp, err := client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	log.Printf("✅ Found %d users (total: %d)", len(listResp.Users), listResp.Total)
	for _, user := range listResp.Users {
		log.Printf("   - %s (%s) [ID: %d]", user.Name, user.Email, user.Id)
	}

	// 4. Test Server Streaming RPC - Watch Users
	log.Println("\n4️⃣ Testing Server Streaming RPC - WatchUsers")
	watchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stream, err := client.WatchUsers(watchCtx, &pb.WatchUsersRequest{})
	if err != nil {
		log.Fatalf("WatchUsers failed: %v", err)
	}

	log.Println("👀 Watching for user updates (10 seconds)...")

	// Start a goroutine to create users while watching
	go func() {
		time.Sleep(2 * time.Second)
		client.CreateUser(ctx, &pb.CreateUserRequest{
			Name:  "Charlie Brown",
			Email: "charlie@example.com",
		})

		time.Sleep(2 * time.Second)
		client.CreateUser(ctx, &pb.CreateUserRequest{
			Name:  "Diana Prince",
			Email: "diana@example.com",
		})
	}()

	watchCount := 0
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("WatchUsers error: %v", err)
			break
		}

		watchCount++
		log.Printf("📺 [%d] %s: %s (%s)", watchCount, resp.Message, resp.User.Name, resp.User.Email)

		// Limit output for demo
		if watchCount >= 10 {
			break
		}
	}

	// 5. Test Client Streaming RPC - Batch Create Users
	log.Println("\n5️⃣ Testing Client Streaming RPC - BatchCreateUsers")
	batchStream, err := client.BatchCreateUsers(ctx)
	if err != nil {
		log.Fatalf("BatchCreateUsers failed: %v", err)
	}

	// Send multiple users
	users := []*pb.CreateUserRequest{
		{Name: "Eve Adams", Email: "eve@example.com"},
		{Name: "Frank Miller", Email: "frank@example.com"},
		{Name: "Grace Lee", Email: "grace@example.com"},
		{Name: "", Email: "invalid@example.com"}, // Invalid user
		{Name: "Henry Wilson", Email: "henry@example.com"},
	}

	log.Printf("📦 Sending %d users in batch...", len(users))
	for i, user := range users {
		if err := batchStream.Send(user); err != nil {
			log.Fatalf("Failed to send user %d: %v", i, err)
		}
		log.Printf("   📤 Sent: %s (%s)", user.Name, user.Email)
		time.Sleep(500 * time.Millisecond) // Small delay for demo
	}

	batchResp, err := batchStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("BatchCreateUsers close failed: %v", err)
	}

	log.Printf("✅ Batch complete: %d users created", batchResp.CreatedCount)
	if len(batchResp.Errors) > 0 {
		log.Printf("⚠️  Errors: %v", batchResp.Errors)
	}

	// Final list to see all users
	log.Println("\n6️⃣ Final user list:")
	finalList, err := client.ListUsers(ctx, &pb.ListUsersRequest{Page: 1, Limit: 20})
	if err != nil {
		log.Fatalf("Final ListUsers failed: %v", err)
	}

	log.Printf("📋 Total users in system: %d", finalList.Total)
	for i, user := range finalList.Users {
		log.Printf("   %d. %s (%s) [ID: %d, Created: %s]",
			i+1, user.Name, user.Email, user.Id,
			time.Unix(user.CreatedAt, 0).Format("15:04:05"))
	}

	log.Println("\n🎉 gRPC client demo completed successfully!")
}
