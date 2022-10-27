package models

import (
	"database/sql"
	"errors"
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
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	var stmt string

	if expires == 1 {
		stmt = `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, now(), now() + interval '1' day) RETURNING id`
	} else if expires == 7 {
		stmt = `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, now(), now() + interval '7' day) RETURNING id`
	} else if expires == 365 {
		stmt = `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, now(), now() + interval '365' day) RETURNING id`
	}

	var id int
	err := m.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > now() AND id = $1`

	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
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

	rows, err := m.DB.Query(stmt)
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
