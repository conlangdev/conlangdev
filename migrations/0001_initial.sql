CREATE TABLE conlangdev_db_status (
    status_key VARCHAR(128) NOT NULL,
    status_value INT NOT NULL,
    PRIMARY KEY(status_key)
);

INSERT INTO conlangdev_db_status (status_key, status_value) VALUES (
    "current_migration", 1
);

CREATE TABLE users (
    id INTEGER NOT NULL AUTOINCREMENT,
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
    id INTEGER NOT NULL AUTOINCREMENT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    endonym VARCHAR(255),
    user_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT uc_user_slug UNIQUE(slug, user_id),
    FOREIGN KEY(user_id) REFERENCES user(id)
);