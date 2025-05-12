CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    role_id INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id)
);

INSERT INTO roles (name , description , level)
VALUES (
    'user' ,
    'A user can create posts and comments',
    1
);
INSERT INTO roles (name , description , level)
VALUES (
    'moderator' ,
    'A moderator can update others posts',
    2
);
INSERT INTO roles (name , description , level)
VALUES (
    'admin' ,
    'An admin can update and delete other users posts',
    3
);