package sql

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	db *sql.DB

	Host     string
	User     string
	Password string
	Database string
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func NewDB(host string, user string, password string, database string) *DB {
	return &DB{
		Host:     host,
		User:     user,
		Password: password,
		Database: database,
	}
}

func (db *DB) Migrate() error {
	log.Info("✨ Beginning migrations")
	if _, err := db.db.Exec("CREATE TABLE IF NOT EXISTS conlangdev_migrations (name VARCHAR(255) PRIMARY KEY);"); err != nil {
		return fmt.Errorf("could not create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationFiles, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error on %s: %w", name, err)
		}
	}

	log.Info("😎 Migrations complete!")
	return nil
}

func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err := tx.QueryRow("SELECT COUNT(*) FROM conlangdev_migrations WHERE name = ?", name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		log.Infof("⏭️ Skipping migration %s.", name)
		return nil
	}

	log.Infof("⏩ Applying migration %s...", name)

	var buffer []byte
	if buffer, err = fs.ReadFile(migrationFiles, name); err != nil {
		return err
	}

	for _, query := range strings.Split(string(buffer), ";\n") {
		if _, err := tx.Exec(query); err != nil {
			return err
		}
	}

	if _, err := tx.Exec("INSERT INTO conlangdev_migrations (name) VALUES (?)", name); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) TestConnection() error {
	log.Info("🍉 Testing database connection")
	var version string
	if err := db.db.QueryRow("SELECT VERSION();").Scan(&version); err != nil {
		return err
	} else if version == "" {
		return errors.New("nothing returned from VERSION() connection test")
	}
	log.Info("🍓 Looking good! Database is healthy.")
	return nil
}

func (db *DB) Open() (err error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", db.User, db.Password, db.Host, db.Database)
	log.Info("📪 Connecting to the database...")
	if db.db, err = sql.Open("mysql", connString); err != nil {
		return err
	}
	log.Infof("📬 Connected as %s@%s on %s!", db.User, db.Host, db.Database)

	if err := db.TestConnection(); err != nil {
		return err
	}

	return db.Migrate()
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
