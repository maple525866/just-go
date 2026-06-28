package sqlcrud

import (
	"database/sql"
	"errors"
)

// Article is the row shape used by the database/sql examples.
type Article struct {
	ID     int64
	Title  string
	Body   string
	Author string
}

// Repository demonstrates a small database/sql repository.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a repository backed by db.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts one article using bound parameters.
func (r *Repository) Create(title, body, author string) (Article, error) {
	result, err := r.db.Exec(`INSERT INTO articles (title, body, author) VALUES (?, ?, ?)`, title, body, author)
	if err != nil {
		return Article{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Article{}, err
	}
	return Article{ID: id, Title: title, Body: body, Author: author}, nil
}

// Get returns one article by ID.
func (r *Repository) Get(id int64) (Article, error) {
	row := r.db.QueryRow(`SELECT id, title, body, author FROM articles WHERE id = ?`, id)
	return scanArticle(row)
}

// Update changes all editable fields and returns the updated article.
func (r *Repository) Update(id int64, title, body, author string) (Article, error) {
	result, err := r.db.Exec(`UPDATE articles SET title = ?, body = ?, author = ? WHERE id = ?`, title, body, author, id)
	if err != nil {
		return Article{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Article{}, err
	}
	if affected == 0 {
		return Article{}, sql.ErrNoRows
	}
	return Article{ID: id, Title: title, Body: body, Author: author}, nil
}

// Delete removes one article by ID.
func (r *Repository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM articles WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// List returns all articles ordered by ID.
func (r *Repository) List() ([]Article, error) {
	rows, err := r.db.Query(`SELECT id, title, body, author FROM articles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Body, &article.Author); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

// FindByTitle demonstrates parameter binding for user input.
func (r *Repository) FindByTitle(title string) ([]Article, error) {
	rows, err := r.db.Query(`SELECT id, title, body, author FROM articles WHERE title = ? ORDER BY id`, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Body, &article.Author); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanArticle(row scanner) (Article, error) {
	var article Article
	if err := row.Scan(&article.ID, &article.Title, &article.Body, &article.Author); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Article{}, err
		}
		return Article{}, err
	}
	return article, nil
}
