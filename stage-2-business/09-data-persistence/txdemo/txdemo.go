package txdemo

import (
	"database/sql"
	"errors"
)

// CommitTwoArticles inserts two rows and commits them as one unit.
func CommitTwoArticles(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO articles (title, body, author) VALUES (?, ?, ?)`, "transaction commit", "first insert commits", "tx"); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`INSERT INTO articles (title, body, author) VALUES (?, ?, ?)`, "transaction commit again", "second insert commits", "tx"); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// RollbackOnError writes a row and then rolls it back after a simulated error.
func RollbackOnError(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO articles (title, body, author) VALUES (?, ?, ?)`, "transaction rollback", "this insert should disappear", "tx"); err != nil {
		_ = tx.Rollback()
		return err
	}
	simulated := errors.New("simulated transaction failure")
	if err := tx.Rollback(); err != nil {
		return err
	}
	return simulated
}
