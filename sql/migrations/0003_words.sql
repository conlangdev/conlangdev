CREATE TABLE words (
    id INT NOT NULL AUTO_INCREMENT,
    uid BIGINT UNSIGNED NOT NULL DEFAULT (UUID_SHORT()),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    headword VARCHAR(255) NOT NULL,
    part_of_speech VARCHAR(255) NOT NULL,
    definition TEXT NOT NULL,
    pronunciation VARCHAR(255),
    grammar_class VARCHAR(255),
    gender VARCHAR(255),
    etymology TEXT,
    notes TEXT,
    language_id INT NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT uc_uid UNIQUE(uid),
    FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE
);