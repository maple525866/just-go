package sqlcrud

import (
	"database/sql"
	"testing"

	"just-go/stage-2-business/09-data-persistence/dbx"
)

func TestRepositoryCRUD(t *testing.T) {
	db := openMigratedDB(t)
	defer db.Close()
	repo := NewRepository(db)

	created, err := repo.Create("SQL basics", "database/sql uses rows and statements", "gopher")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("created article should have ID")
	}

	got, err := repo.Get(created.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Title != "SQL basics" {
		t.Fatalf("title = %q, want SQL basics", got.Title)
	}

	updated, err := repo.Update(created.ID, "SQL updated", "updates use bound parameters", "gopher")
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Title != "SQL updated" {
		t.Fatalf("updated title = %q", updated.Title)
	}

	items, err := repo.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("list length = %d, want 1", len(items))
	}

	if err := repo.Delete(created.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := repo.Get(created.ID); err == nil {
		t.Fatal("Get after Delete succeeded, want error")
	}
}

func TestFindByTitleUsesBoundParameters(t *testing.T) {
	db := openMigratedDB(t)
	defer db.Close()
	repo := NewRepository(db)
	if _, err := repo.Create("safe title", "body", "gopher"); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	items, err := repo.FindByTitle("safe title' OR 1=1 --")
	if err != nil {
		t.Fatalf("FindByTitle returned error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("SQL injection-like title returned %d rows, want 0", len(items))
	}

	if _, err := repo.Create("safe title' OR 1=1 --", "literal input", "attacker"); err != nil {
		t.Fatalf("Create literal title returned error: %v", err)
	}
	items, err = repo.FindByTitle("safe title' OR 1=1 --")
	if err != nil {
		t.Fatalf("FindByTitle literal returned error: %v", err)
	}
	if len(items) != 1 || items[0].Author != "attacker" {
		t.Fatalf("literal query result = %+v, want attacker row", items)
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
