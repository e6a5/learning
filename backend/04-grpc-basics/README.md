# ⚡ 04 - gRPC Basics

This project demonstrates how to build gRPC services in Go with both unary and streaming RPCs. Learn the fundamentals of Protocol Buffers, gRPC communication patterns, and containerized microservices using modern buf.build toolchain.

---

## 🎯 What You'll Learn

- Define services and messages with Protocol Buffers
- Use buf.build for modern protobuf development workflow
- Implement gRPC servers with unary and streaming RPCs
- Build gRPC clients for different communication patterns
- Use Docker for containerized gRPC services
- Test gRPC services with grpcurl and custom clients
- Understand the four types of gRPC communication

---

## 🧱 Stack

- **Golang** - Backend implementation
- **gRPC** - High-performance RPC framework
- **Protocol Buffers** - Interface definition language
- **buf.build** - Modern protobuf toolchain
- **Docker & Docker Compose** - Containerization
- **grpcurl** - Command-line gRPC client (for testing)

---

## 🚀 Quick Start

```bash
# Start the gRPC server
make up

# Check if server is running
make ps

# Run the full demo (server + client)
make demo
```

## ⚡ gRPC Service Methods

| Method | Type | Description | Use Case |
|--------|------|-------------|----------|
| `CreateUser` | Unary | Create a single user | Simple CRUD operations |
| `GetUser` | Unary | Get user by ID | Data retrieval |
| `ListUsers` | Unary | List users with pagination | Bulk data queries |
| `WatchUsers` | Server Streaming | Real-time user updates | Live notifications |
| `BatchCreateUsers` | Client Streaming | Bulk user creation | Batch operations |

---

## 🔄 RPC Communication Patterns

### 1️⃣ Unary RPC (Request-Response)
```bash
# Create a user
make test-create

# Get a user
make test-get

# List users
make test-list
```

### 2️⃣ Server Streaming RPC (Real-time Updates)
```bash
# Watch for user updates
make test-watch
```

### 3️⃣ Client Streaming RPC (Batch Operations)
The client sends multiple messages to the server and gets a single response.

---

## 🧪 Testing the Service

### Using the Built-in Client Demo
```bash
# Run complete demo with all RPC patterns
make demo

# View demo output
make logs-client
```

### Using grpcurl (Command Line)
```bash
# Create a user
grpcurl -plaintext -d '{"name":"Alice","email":"alice@test.com"}' \
  localhost:50051 user.UserService/CreateUser

# Get user by ID
grpcurl -plaintext -d '{"id":1}' \
  localhost:50051 user.UserService/GetUser

# List all users
grpcurl -plaintext -d '{"page":1,"limit":5}' \
  localhost:50051 user.UserService/ListUsers

# Watch users in real-time (Ctrl+C to stop)
grpcurl -plaintext -d '{}' \
  localhost:50051 user.UserService/WatchUsers
```

### Service Discovery
```bash
# List available services
make list-services

# Describe the UserService
make describe-service
```

---

## 🔧 Development Commands

### Docker & Services
```bash
# Build and start server
make up

# Run full demo (server + client)
make demo

# View server logs
make logs

# View client logs
make logs-client

# Stop services
make down

# Rebuild everything
make rebuild
```

### Buf & Protobuf Workflow
```bash
# Setup buf development environment
make dev-setup

# Generate protobuf code
make gen

# Lint protobuf files
make lint

# Format protobuf files
make format

# Check for breaking changes
make breaking

# Run full protobuf workflow
make proto-workflow

# See all commands
make help
```

---

## 📋 Buf Configuration

### buf.yaml (Project Configuration)
```yaml
version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
deps:
  - buf.build/googleapis/googleapis
```

### buf.gen.yaml (Code Generation)
```yaml
version: v1
managed:
  enabled: true
  go_package_prefix:
    default: github.com/e6a5/learning/backend/04-grpc-basics
plugins:
  - plugin: buf.build/protocolbuffers/go
    out: .
    opt:
      - paths=source_relative
  - plugin: buf.build/grpc/go
    out: .
    opt:
      - paths=source_relative
```

---

## 📋 Protocol Buffer Definition

```protobuf
service UserService {
  // Unary RPCs
  rpc CreateUser(CreateUserRequest) returns (UserResponse);
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  
  // Streaming RPCs
  rpc WatchUsers(WatchUsersRequest) returns (stream UserResponse);
  rpc BatchCreateUsers(stream CreateUserRequest) returns (BatchCreateResponse);
}

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  int64 created_at = 4;
}
```

---

## 🔍 buf.build vs protoc

| Feature | buf.build | protoc |
|---------|-----------|--------|
| **Setup** | Single binary | Multiple tools & plugins |
| **Linting** | Built-in comprehensive linting | Manual setup required |
| **Breaking Changes** | Automatic detection | Manual tracking |
| **Code Generation** | Centralized configuration | Command-line arguments |
| **Dependencies** | Built-in dependency management | Manual management |
| **Performance** | Optimized for speed | Standard performance |

---

## 📦 File Overview

| File | Purpose |
|------|---------|
| `buf.yaml` | Buf project configuration |
| `buf.gen.yaml` | Code generation configuration |
| `proto/user.proto` | Protocol Buffer service definition |
| `server/main.go` | gRPC server implementation |
| `client/main.go` | gRPC client demo |
| `Dockerfile.server` | Server container build |
| `Dockerfile.client` | Client container build |
| `compose.yml` | Multi-service orchestration |
| `Makefile` | Development and testing commands |

---

## 🐳 Docker Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │────│   gRPC Server   │
│   (Demo)        │    │   :50051        │
│                 │    │   (Protocol     │
│                 │    │    Buffers)     │
└─────────────────┘    └─────────────────┘
```

- **Server Container**: Runs the gRPC server on port 50051
- **Client Container**: Demonstrates all RPC patterns
- **Network**: Isolated Docker network for gRPC communication
- **Health Checks**: Ensures server is ready before client starts
- **Buf Integration**: Each container generates protobuf code using buf

---

## 🔧 Code Generation with Buf

Buf generates Go code from `.proto` files using the configuration in `buf.gen.yaml`:

```bash
# Generate Go code from proto files
buf generate

# Generate with specific configuration
buf generate --template buf.gen.yaml
```

Generated files:
- `proto/user.pb.go` - Message types
- `proto/user_grpc.pb.go` - Service interfaces

### Benefits of Buf:
- **Faster builds** - Optimized compilation
- **Built-in linting** - Consistent code style
- **Breaking change detection** - API evolution safety
- **Dependency management** - Handle external dependencies
- **Remote packages** - Use buf.build registry

---

## 🚨 Common Issues & Solutions

### Server Not Starting
```bash
# Check server logs
make logs

# Verify port is not in use
netstat -an | grep 50051
```

### Client Connection Issues
```bash
# Ensure server is healthy
make ps

# Check network connectivity
docker compose exec client ping server
```

### Buf Setup Issues
```bash
# Install buf locally
make dev-setup

# Or manually:
go install github.com/bufbuild/buf/cmd/buf@latest

# Verify installation
buf --version
```

### Proto Generation Issues
```bash
# Clean and regenerate
make clean
make gen

# Check buf configuration
buf config ls-files
buf config ls-breaking-rules
buf config ls-lint-rules
```

---

## 🎭 Demo Scenarios

### Scenario 1: Basic CRUD Operations
1. Start server: `make up`
2. Create users: `make test-create`
3. Retrieve users: `make test-get`
4. List all users: `make test-list`

### Scenario 2: Real-time Streaming
1. Start watching: `make test-watch` (in one terminal)
2. Create users: `make test-create` (in another terminal)
3. Observe real-time updates

### Scenario 3: Full Demo
1. Run complete demo: `make demo`
2. Watch the client demonstrate all RPC patterns
3. See batch operations and streaming in action

### Scenario 4: Development Workflow
1. Setup environment: `make dev-setup`
2. Run protobuf workflow: `make proto-workflow`
3. Make changes to proto files
4. Test and iterate: `make lint && make gen && make up`

---

## 💡 Next Steps

Try these experiments:
1. **Add New Services**: Extend the proto with more RPC methods
2. **Buf Registry**: Publish your protos to buf.build registry
3. **Authentication**: Implement gRPC interceptors for auth
4. **Error Handling**: Add custom error codes and details
5. **Bidirectional Streaming**: Implement chat-like functionality
6. **Load Balancing**: Set up multiple server instances
7. **TLS Security**: Add encrypted gRPC communication
8. **Metrics**: Implement gRPC metrics and monitoring
9. **Gateway**: Add REST-to-gRPC gateway with grpc-gateway
10. **Buf Breaking**: Set up CI/CD with breaking change detection 