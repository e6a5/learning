-- Error Handling Lab Database Schema
-- This script initializes the database with tables for testing error handling scenarios

-- Users table for basic CRUD operations
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_joined_at (joined_at)
);

-- Error logs table for tracking application errors
CREATE TABLE IF NOT EXISTS error_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    error_type VARCHAR(50) NOT NULL,
    error_code VARCHAR(100) NOT NULL,
    error_message TEXT NOT NULL,
    request_id VARCHAR(100),
    endpoint VARCHAR(255),
    user_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_error_type (error_type),
    INDEX idx_created_at (created_at),
    INDEX idx_request_id (request_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Circuit breaker events table for monitoring
CREATE TABLE IF NOT EXISTS circuit_breaker_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    event_type ENUM('opened', 'closed', 'half_open') NOT NULL,
    failure_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_service_name (service_name),
    INDEX idx_created_at (created_at)
);

-- Insert sample data for testing
INSERT INTO users (name, email) VALUES 
    ('Alice Johnson', 'alice@example.com'),
    ('Bob Smith', 'bob@example.com'),
    ('Carol Wilson', 'carol@example.com')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- Create a view for error statistics
CREATE OR REPLACE VIEW error_statistics AS
SELECT 
    error_type,
    error_code,
    COUNT(*) as occurrence_count,
    MAX(created_at) as last_occurrence,
    DATE(created_at) as error_date
FROM error_logs 
GROUP BY error_type, error_code, DATE(created_at)
ORDER BY occurrence_count DESC; 