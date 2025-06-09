package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomMetric_Validate(t *testing.T) {
	tests := []struct {
		name    string
		metric  CustomMetric
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid counter metric",
			metric: CustomMetric{
				Name:  "requests_total",
				Type:  "counter",
				Value: 42.0,
				Labels: map[string]string{
					"method": "GET",
					"status": "200",
				},
			},
			wantErr: false,
		},
		{
			name: "valid gauge metric",
			metric: CustomMetric{
				Name:  "memory_usage",
				Type:  "gauge",
				Value: 1024.5,
			},
			wantErr: false,
		},
		{
			name: "valid histogram metric",
			metric: CustomMetric{
				Name:  "response_time",
				Type:  "histogram",
				Value: 0.125,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			metric: CustomMetric{
				Name:  "",
				Type:  "counter",
				Value: 1.0,
			},
			wantErr: true,
			errMsg:  "Metric name is required",
		},
		{
			name: "name too long",
			metric: CustomMetric{
				Name:  "this_is_a_very_long_metric_name_that_exceeds_the_maximum_allowed_length_of_one_hundred_characters_limit",
				Type:  "counter",
				Value: 1.0,
			},
			wantErr: true,
			errMsg:  "Metric name must be less than 100 characters",
		},
		{
			name: "empty type",
			metric: CustomMetric{
				Name:  "valid_name",
				Type:  "",
				Value: 1.0,
			},
			wantErr: true,
			errMsg:  "Metric type is required",
		},
		{
			name: "invalid type",
			metric: CustomMetric{
				Name:  "valid_name",
				Type:  "invalid_type",
				Value: 1.0,
			},
			wantErr: true,
			errMsg:  "Metric type must be counter, gauge, or histogram",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metric.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewHealthCheck(t *testing.T) {
	tests := []struct {
		name     string
		reqName  string
		message  string
		status   HealthStatus
		duration time.Duration
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid health check",
			reqName:  "database",
			message:  "Connection successful",
			status:   HealthStatusHealthy,
			duration: 50 * time.Millisecond,
			wantErr:  false,
		},
		{
			name:     "degraded service",
			reqName:  "external_api",
			message:  "Slow response",
			status:   HealthStatusDegraded,
			duration: 2 * time.Second,
			wantErr:  false,
		},
		{
			name:     "unhealthy service",
			reqName:  "cache",
			message:  "Connection failed",
			status:   HealthStatusUnhealthy,
			duration: 5 * time.Second,
			wantErr:  false,
		},
		{
			name:     "empty name",
			reqName:  "",
			message:  "Some message",
			status:   HealthStatusHealthy,
			duration: 100 * time.Millisecond,
			wantErr:  true,
			errMsg:   "Health check name is required",
		},
		{
			name:     "name too long",
			reqName:  "this_is_a_very_long_health_check_name_that_exceeds_limit",
			message:  "Some message",
			status:   HealthStatusHealthy,
			duration: 100 * time.Millisecond,
			wantErr:  true,
			errMsg:   "Health check name must be less than 50 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check, err := NewHealthCheck(tt.reqName, tt.message, tt.status, tt.duration)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, check)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, check)
				assert.Equal(t, tt.reqName, check.Name)
				assert.Equal(t, tt.message, check.Message)
				assert.Equal(t, tt.status, check.Status)
				assert.Equal(t, tt.duration, check.Duration)
				assert.False(t, check.Timestamp.IsZero())
			}
		})
	}
}

func TestHealthResponse_IsHealthy(t *testing.T) {
	tests := []struct {
		name     string
		response HealthResponse
		want     bool
	}{
		{
			name: "healthy status",
			response: HealthResponse{
				Status: HealthStatusHealthy,
			},
			want: true,
		},
		{
			name: "degraded status",
			response: HealthResponse{
				Status: HealthStatusDegraded,
			},
			want: false,
		},
		{
			name: "unhealthy status",
			response: HealthResponse{
				Status: HealthStatusUnhealthy,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.IsHealthy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHealthResponse_HasCriticalFailures(t *testing.T) {
	tests := []struct {
		name     string
		response HealthResponse
		want     bool
	}{
		{
			name: "no critical failures",
			response: HealthResponse{
				Checks: []HealthCheck{
					{Status: HealthStatusHealthy},
					{Status: HealthStatusDegraded},
				},
			},
			want: false,
		},
		{
			name: "has critical failure",
			response: HealthResponse{
				Checks: []HealthCheck{
					{Status: HealthStatusHealthy},
					{Status: HealthStatusUnhealthy},
				},
			},
			want: true,
		},
		{
			name: "all unhealthy",
			response: HealthResponse{
				Checks: []HealthCheck{
					{Status: HealthStatusUnhealthy},
					{Status: HealthStatusUnhealthy},
				},
			},
			want: true,
		},
		{
			name: "no checks",
			response: HealthResponse{
				Checks: []HealthCheck{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.HasCriticalFailures()
			assert.Equal(t, tt.want, got)
		})
	}
}
