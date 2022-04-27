CREATE TABLE words (
    id INT NOT NULL AUTO_INCREMENT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    headword VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    definition TEXT NOT NULL,
    etymology TEXT,
    notes TEXT,
    language_id INT NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT uc_language_slug UNIQUE(slug, language_id),
    FOREIGN KEY (language_id) REFERENCES languages(id)
);