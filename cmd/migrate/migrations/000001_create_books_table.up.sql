CREATE TABLE IF NOT EXISTS books (
    id BIGSERIAL  PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(50) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    stock INT NOT NULL,
    tags TEXT[],
    pages INT NOT NULL,
    cover_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);