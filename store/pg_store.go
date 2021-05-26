package store

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type pgStore struct {
	db *sql.DB
}

func NewPostgresStore(connectionStr string) (*pgStore, error) {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %v: %v", connectionStr, err)
	}
	return &pgStore{db: db}, nil
}

func (s *pgStore) Ping() error {
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("couldn't connect db: %v", err)
	}
	return nil
}

func (s *pgStore) Create(q Quote) (Quote, error) {
	stmt := `
INSERT INTO quotes(author, text)
VALUES($1, $2)
RETURNING id
`

	var id int
	if err := s.db.QueryRow(stmt, q.Author, q.Text).Scan(&id); err != nil {
		return Quote{}, err
	}

	q.ID = id

	return q, nil
}

func (s *pgStore) Get(id int) (Quote, error) {
	stmt := `
SELECT id, author, text
FROM quotes
WHERE id = $1
`
	var q Quote
	if err := s.db.QueryRow(stmt, id).Scan(&q.ID, &q.Author, &q.Text); err != nil {
		return Quote{}, err
	}

	return q, nil
}

func (s *pgStore) GetAll() ([]Quote, error) {
	rows, err := s.db.Query("SELECT id, author, text FROM quotes")
	if err != nil {
		return []Quote{}, err
	}

	defer rows.Close()

	quotes := []Quote{}
	for rows.Next() {
		var q Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Text); err != nil {
			return []Quote{}, err
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (s *pgStore) findMaxID() (int, error) {
	var maxID int
	if err := s.db.QueryRow("SELECT MAX(id) FROM quotes").Scan(&maxID); err != nil {
		return -1, err
	}

	return maxID, nil
}

func (s *pgStore) GetRandom() (Quote, error) {
	maxID, err := s.findMaxID()
	if err != nil {
		return Quote{}, err
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(maxID-1) + 1
	return s.Get(id)
}

func (s *pgStore) Clean() error {
	if _, err := s.db.Exec("DELETE FROM quotes"); err != nil {
		return err
	}
	return nil
}
