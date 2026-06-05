CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(50) PRIMARY KEY
);

INSERT INTO roles (name) VALUES ('admin'), ('user')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user'
    REFERENCES roles(name);

-- Give the seeded Alice an admin role so there is an admin to test with.
UPDATE users SET role = 'admin' WHERE email = 'alice@example.com';
