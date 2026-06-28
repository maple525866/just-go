package dbx

import (
	"database/sql"
	"strings"
	"testing"
	"time"
)

func TestApplyMigrationCreatesArticlesTable(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := ApplyMigration(db, "../migrations/001_create_articles.sql"); err != nil {
		t.Fatalf("ApplyMigration returned error: %v", err)
	}

	_, err := db.Exec(`INSERT INTO articles (title, body, author) VALUES (?, ?, ?)`, "SQL basics", "Rows are persisted in tables", "gopher")
	if err != nil {
		t.Fatalf("insert after migration failed: %v", err)
	}
}

func TestPoolConfigSummary(t *testing.T) {
	config := PoolConfig{MaxOpenConns: 4, MaxIdleConns: 2, ConnMaxLifetime: 30 * time.Minute}

	summary := config.Summary()

	for _, want := range []string{"max open=4", "max idle=2", "lifetime=30m0s"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary %q does not contain %q", summary, want)
		}
	}
}

func TestConfigurePoolAppliesSettings(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	config := PoolConfig{MaxOpenConns: 3, MaxIdleConns: 1, ConnMaxLifetime: time.Minute}

	ConfigurePool(db, config)
	stats := db.Stats()

	if stats.MaxOpenConnections != 3 {
		t.Fatalf("MaxOpenConnections = %d, want 3", stats.MaxOpenConnections)
	}
}

func TestOpenMemorySharesSchemaAcrossConnections(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	ConfigurePool(db, PoolConfig{MaxOpenConns: 2, MaxIdleConns: 2, ConnMaxLifetime: time.Minute})
	if err := ApplyMigration(db, "../migrations/001_create_articles.sql"); err != nil {
		t.Fatalf("ApplyMigration returned error: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Begin returned error: %v", err)
	}
	defer tx.Rollback()

	_, err = db.Exec(`INSERT INTO articles (title, body, author) VALUES (?, ?, ?)`, "shared schema", "second connection sees migration", "gopher")
	if err != nil {
		t.Fatalf("insert through second connection failed: %v", err)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory returned error: %v", err)
	}
	return db
}
