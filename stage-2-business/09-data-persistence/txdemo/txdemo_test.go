package txdemo

import (
	"database/sql"
	"testing"

	"just-go/stage-2-business/09-data-persistence/dbx"
	"just-go/stage-2-business/09-data-persistence/sqlcrud"
)

func TestCommitTwoArticles(t *testing.T) {
	db := openMigratedDB(t)
	defer db.Close()

	if err := CommitTwoArticles(db); err != nil {
		t.Fatalf("CommitTwoArticles returned error: %v", err)
	}

	items, err := sqlcrud.NewRepository(db).List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("article count = %d, want 2", len(items))
	}
}

func TestRollbackOnError(t *testing.T) {
	db := openMigratedDB(t)
	defer db.Close()

	if err := RollbackOnError(db); err == nil {
		t.Fatal("RollbackOnError returned nil, want error")
	}

	items, err := sqlcrud.NewRepository(db).List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("article count after rollback = %d, want 0", len(items))
	}
}

func openMigratedDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := dbx.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory returned error: %v", err)
	}
	if err := dbx.ApplyMigration(db, "../migrations/001_create_articles.sql"); err != nil {
		t.Fatalf("ApplyMigration returned error: %v", err)
	}
	return db
}
