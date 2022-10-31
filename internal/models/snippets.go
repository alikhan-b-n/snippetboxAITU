package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	timeExp := time.Now().AddDate(0, 0, expires)

	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, now(), $3) RETURNING id`

	var id int
	//err := m.DB.QueryRow(stmt, title, content, timeExp).Scan(&id)
	err := m.DB.QueryRow(context.Background(), stmt, title, content, timeExp).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > now() AND id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > now() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	snippets := []*Snippet{}

	for rows.Next() {

		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
