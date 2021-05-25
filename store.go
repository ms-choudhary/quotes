package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init(module, connectionStr string) error {
	var err error
	db, err = sql.Open(module, connectionStr)
	if err != nil {
		return fmt.Errorf("cannot open db: (%v) %v: %v", module, connectionStr, err)
	}
	return nil
}

func Create(q Quote) (Quote, error) {
	stmt := `
INSERT INTO quotes(author, text)
VALUES($1, $2)
RETURNING id
`

	var id int
	if err := db.QueryRow(stmt, q.Author, q.Text).Scan(&id); err != nil {
		return Quote{}, err
	}

	q.ID = id

	return q, nil
}

func Get(id int) (Quote, error) {
	stmt := `
SELECT id, author, text
FROM quotes
WHERE id = $1
`
	var q Quote
	if err := db.QueryRow(stmt, id).Scan(&q.ID, &q.Author, &q.Text); err != nil {
		return Quote{}, err
	}

	return q, nil
}

func GetAll() ([]Quote, error) {
	rows, err := db.Query("SELECT id, author, text FROM quotes")
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

func findMaxID() (int, error) {
	var maxID int
	if err := db.QueryRow("SELECT MAX(id) FROM quotes").Scan(&maxID); err != nil {
		return -1, err
	}

	return maxID, nil
}

func GetRandom() (Quote, error) {
	maxID, err := findMaxID()
	if err != nil {
		return Quote{}, err
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(maxID-1) + 1
	return Get(id)
}
