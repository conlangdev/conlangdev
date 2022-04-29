CREATE TABLE users (
    id INTEGER NOT NULL AUTO_INCREMENT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT uc_username UNIQUE(username),
    CONSTRAINT uc_email UNIQUE(email)
);