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

CREATE TABLE languages (
    id INTEGER NOT NULL AUTO_INCREMENT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    endonym VARCHAR(255),
    user_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT uc_user_slug UNIQUE(slug, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);