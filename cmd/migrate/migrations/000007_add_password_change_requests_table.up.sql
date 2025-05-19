CREATE TABLE IF NOT EXISTS password_change_requests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ,
    token VARCHAR(64) NOT NULL UNIQUE,
    expiry TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);