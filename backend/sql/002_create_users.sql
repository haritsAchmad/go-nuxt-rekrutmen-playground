CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL,
    password_hash VARCHAR(64) NOT NULL,
    password_salt VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'aktif',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (LOWER(email));

INSERT INTO users (name, email, role, password_hash, password_salt, status)
VALUES (
    'Admin Demo',
    'admin@example.com',
    'admin',
    '9647d1b5a75928bd45b120230ea9e032f10538ed5948da78b593bde75e5991ce',
    'demo_salt',
    'aktif'
)
ON CONFLICT (email) DO NOTHING;

-- Demo login:
-- email: admin@example.com
-- password: admin123
