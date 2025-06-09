# 📊 08-monitoring: System Observability & Health Monitoring

**Learning Question**: *"How do you know if your system is healthy?"*

This module explores comprehensive **monitoring and observability** patterns for production backend systems. Learn how to track system health, collect metrics, monitor performance, and detect issues before they impact users.

---

## 🎯 Learning Objectives

By the end of this module, you'll understand:

- **Health Checks**: Liveness, readiness, and dependency health monitoring
- **Metrics Collection**: Custom metrics, request/response tracking, system metrics
- **Performance Monitoring**: Response times, throughput, resource usage
- **Alerting Patterns**: How to detect and respond to system issues
- **Observability Tools**: Prometheus, Grafana, structured logging
- **Production Monitoring**: Real-world monitoring strategies

---

## 🏗️ Architecture Overview

```
08-monitoring/                     # Observability module
├── main.go                       # Service orchestration (72 lines)
├── internal/
│   ├── models/                   # Metric data structures (123 lines)
│   │   └── metrics.go           # Health checks, metrics, validation
│   ├── repository/               # Metrics storage & health checks (201 lines)
│   │   └── metrics.go           # In-memory metrics, health checkers
│   ├── handlers/                 # HTTP monitoring endpoints (273 lines)
│   │   └── monitoring.go        # Health, metrics, status endpoints
│   ├── middleware/               # Request monitoring middleware (142 lines)
│   │   └── monitoring.go        # Metrics collection, logging
│   └── utils/                    # Response utilities (21 lines)
│       └── response.go          # JSON response helpers
├── compose.yml                   # Full monitoring stack
├── README.md                     # This documentation
└── UNIX_PHILOSOPHY.md           # Architecture principles
```

**Total**: 832 lines across 7 focused files (avg: 119 lines/file)

---

## 🚀 Quick Start

### 1. Start the Monitoring Stack
```bash
# Start full monitoring infrastructure
docker compose up -d

# Or start just the application
go run main.go
```

### 2. Explore Health Endpoints
```bash
# Comprehensive health check
curl http://localhost:8080/health

# Kubernetes-style liveness probe
curl http://localhost:8080/health/live

# Kubernetes-style readiness probe
curl http://localhost:8080/health/ready
```

### 3. View Metrics
```bash
# Prometheus-style metrics
curl http://localhost:8080/metrics

# Custom JSON metrics
curl http://localhost:8080/api/metrics

# System information
curl http://localhost:8080/api/system

# Application status overview
curl http://localhost:8080/api/status
```

### 4. Generate Test Metrics
```bash
# Normal request
curl http://localhost:8080/api/demo

# Simulate errors
curl http://localhost:8080/api/demo?error=500

# Simulate slow response
curl http://localhost:8080/api/demo?delay=2000
```

---

## 📊 Monitoring Features

### 🟢 Health Checks

**Comprehensive Health**: `/health`
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 3600,
  "environment": "development",
  "checks": [
    {
      "name": "database",
      "status": "healthy",
      "duration_ms": 45,
      "details": {"type": "database"}
    },
    {
      "name": "api",
      "status": "healthy", 
      "duration_ms": 123,
      "details": {"type": "external_service"}
    }
  ]
}
```

**Liveness Probe**: `/health/live`
- Always returns 200 if process is running
- Used by Kubernetes for restart decisions

**Readiness Probe**: `/health/ready`
- Returns 503 if dependencies are unhealthy
- Used by load balancers for traffic routing

### 📈 Metrics Collection

**Request Metrics**:
- Request count by endpoint and method
- Response times and status codes
- Request/response sizes
- Client IP and User-Agent tracking

**System Metrics**:
- Memory usage (heap, total)
- Goroutine count
- CPU usage (basic)
- Application uptime

**Custom Metrics**:
```bash
# Submit custom metric
curl -X POST http://localhost:8080/api/metrics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "user_registrations_total",
    "type": "counter",
    "value": 1,
    "labels": {"source": "web", "plan": "premium"}
  }'
```

### 🔍 Observability Dashboard

Access monitoring tools:
- **Application**: http://localhost:8080
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

---

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `VERSION` | `1.0.0` | Application version |
| `ENVIRONMENT` | `development` | Deployment environment |

### Health Check Configuration

Health checkers are configured in `main.go`:
```go
healthCheckers := []repository.HealthChecker{
    repository.NewDatabaseHealthChecker("database", "mysql://localhost:3306"),
    repository.NewExternalServiceHealthChecker("api", "https://httpbin.org/status/200"),
}
```

Add custom health checks by implementing the `HealthChecker` interface:
```go
type HealthChecker interface {
    Check(ctx context.Context) models.HealthCheck
}
```

---

## 🧪 Testing & Validation

### Manual Testing
```bash
# Generate load for testing
for i in {1..100}; do
  curl -s http://localhost:8080/api/demo &
done

# Check metrics after load
curl http://localhost:8080/api/metrics | jq '.request_metrics'

# Test error scenarios
curl http://localhost:8080/api/demo?error=500
curl http://localhost:8080/api/demo?delay=3000
```

### Health Check Testing
```bash
# Test with dependency failures
docker compose stop mysql

# Check health status
curl http://localhost:8080/health

# Restart dependency
docker compose start mysql
```

### Performance Monitoring
```bash
# Monitor system metrics under load
watch -n 1 'curl -s http://localhost:8080/api/system | jq .system_metrics'
```

---

## 📊 Monitoring Best Practices

### 1. **Health Check Patterns**
- ✅ **Liveness**: Simple "am I alive?" check
- ✅ **Readiness**: "Am I ready to serve traffic?"
- ✅ **Dependency Health**: Check external services
- ✅ **Timeout Handling**: Always use timeouts for checks

### 2. **Metrics Strategy**
- ✅ **RED Method**: Rate, Errors, Duration
- ✅ **USE Method**: Utilization, Saturation, Errors
- ✅ **Custom Business Metrics**: Domain-specific measurements
- ✅ **Structured Logging**: Machine-readable logs

### 3. **Performance Monitoring**
- ✅ **Request Tracking**: Every HTTP request
- ✅ **Resource Monitoring**: Memory, CPU, goroutines
- ✅ **Error Tracking**: Categorized error metrics
- ✅ **Alerting Thresholds**: Define clear SLOs

### 4. **Production Readiness**
- ✅ **Graceful Shutdown**: Clean resource cleanup
- ✅ **Circuit Breakers**: Protect against cascading failures
- ✅ **Rate Limiting**: Prevent abuse
- ✅ **Security Headers**: CORS, security headers

---

## 🔍 Real-World Applications

### Kubernetes Integration
```yaml
# Health check configuration
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Prometheus Configuration
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'monitoring-service'
    static_configs:
      - targets: ['monitoring-app:8080']
    scrape_interval: 15s
    metrics_path: '/metrics'
```

### Grafana Dashboards
- Request rate and latency
- Error rate by endpoint
- System resource usage
- Health check status
- Custom business metrics

---

## ⚡ Performance Characteristics

### Metrics Collection Overhead
- **Memory**: ~1-5MB for 10k metrics
- **CPU**: <1% overhead for metric collection
- **Latency**: <1ms added per request
- **Storage**: In-memory with configurable retention

### Health Check Performance
- **Database Check**: ~50-100ms
- **External Service**: ~100-500ms
- **Overall Health**: <1s with timeouts
- **Concurrent Checks**: Parallel execution

---

## 🎓 Learning Outcomes

After working through this module:

1. **Health Monitoring**: You can implement comprehensive health checks
2. **Metrics Collection**: You understand how to collect and expose metrics
3. **Performance Tracking**: You can monitor system performance
4. **Observability**: You know how to make systems observable
5. **Production Monitoring**: You understand real-world monitoring needs

---

## 🔗 Related Concepts

- **07-error-handling**: Graceful error handling and recovery
- **09-config-mgmt**: Environment-based configuration
- **10-performance**: Performance optimization techniques
- **14-deployment**: Production deployment with monitoring

---

## 🚨 Common Pitfalls

1. **Over-monitoring**: Collecting too many low-value metrics
2. **Timeout Issues**: Not setting proper timeouts on health checks
3. **Recursive Monitoring**: Monitoring endpoints triggering more monitoring
4. **Storage Growth**: Unbounded metric storage
5. **Alert Fatigue**: Too many false positive alerts

---

**Next Steps**: Explore `09-config-mgmt` to learn environment-based configuration management, or dive into `10-performance` for optimization techniques.

**Production Tip**: Start with basic health checks and essential metrics, then gradually add more detailed monitoring based on actual operational needs! 