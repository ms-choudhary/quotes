package store

import (
	"testing"
)

func TestCreateGetOne(t *testing.T) {
	db, err := NewPostgresStore("postgres://postgres@localhost/production_test?sslmode=disable")

	q, err := db.Create(Quote{Author: "TestAuth", Text: "Testing"})
	if err != nil {
		t.Fatal(err)
	}

	q, err = db.Get(q.ID)
	if err != nil {
		t.Fatal(err)
	}

	if q.Author != "TestAuth" || q.Text != "Testing" {
		t.Error("get is different from created")
	}
}
