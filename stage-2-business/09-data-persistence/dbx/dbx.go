package dbx

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// PoolConfig captures the connection-pool settings worth learning first.
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// OpenMemory opens a shared SQLite memory database for deterministic local tests.
func OpenMemory() (*sql.DB, error) {
	return sql.Open("sqlite3", "file:just_go_stage09?mode=memory&cache=shared")
}

// ConfigurePool applies connection-pool settings to a database handle.
func ConfigurePool(db *sql.DB, config PoolConfig) {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
}

// Summary returns a short human-readable pool configuration report.
func (c PoolConfig) Summary() string {
	return fmt.Sprintf("max open=%d, max idle=%d, lifetime=%s", c.MaxOpenConns, c.MaxIdleConns, c.ConnMaxLifetime)
}

// ApplyMigration executes a SQL migration file.
func ApplyMigration(db *sql.DB, path string) error {
	migration, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(migration))
	return err
}
