CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT 'system',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) DEFAULT 'system',
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
);

INSERT INTO users (username, email)
SELECT 'Alice', 'alice@example.com'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'alice@example.com');

CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(50) PRIMARY KEY
);

INSERT INTO roles (name) VALUES ('admin'), ('user')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user'
    REFERENCES roles(name);

INSERT INTO users (username, email, role)
SELECT 'John', 'john@example.com', 'user'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'john@example.com');

INSERT INTO users (username, email, role)
SELECT 'Eve', 'eve@example.com', 'admin'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'eve@example.com');

-- Give the seeded Alice an admin role so there is an admin to test with.
UPDATE users SET role = 'admin' WHERE email = 'alice@example.com';
