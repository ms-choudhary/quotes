package main

import (
	"os/exec"
	"testing"
)

func setup(t *testing.T) {
	if db == nil {
		if err := exec.Command("bash", "./hack/migration.sh", "production_test").Run(); err != nil {
			t.Fatal(err)
		}

		err := Init("postgres", "postgres://postgres@localhost/production?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func cleanup(t *testing.T) {
	if err := exec.Command("bash", "./hack/cleanup.sh", "production_test").Run(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateGetOne(t *testing.T) {
	setup(t)

	q, err := Create(Quote{Author: "TestAuth", Text: "Testing"})
	if err != nil {
		t.Fatal(err)
	}

	q, err = Get(q.ID)
	if err != nil {
		t.Fatal(err)
	}

	if q.Author != "TestAuth" || q.Text != "Testing" {
		t.Error("get is different from created")
	}

	//cleanup(t)
}
