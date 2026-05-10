-- Create users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed users with bcrypt hashed passwords
-- All passwords are: password123
-- Hash generated with bcrypt cost 10
INSERT INTO users (username, email, password, created_at, updated_at) VALUES
('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('john_doe', 'john.doe@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('jane_smith', 'jane.smith@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('bob_wilson', 'bob.wilson@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('alice_johnson', 'alice.johnson@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('charlie_brown', 'charlie.brown@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('david_lee', 'david.lee@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('emma_davis', 'emma.davis@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('frank_miller', 'frank.miller@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW()),
('grace_taylor', 'grace.taylor@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhJrAKMQHy5SJE1OG/5i5z3iUxoS', NOW(), NOW());
