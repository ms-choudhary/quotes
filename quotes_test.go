package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ms-choudhary/quotes/store"
)

func TestCreateHandler(t *testing.T) {
	buf := new(bytes.Buffer)
	q := store.Quote{Author: "Test", Text: "Test"}
	if err := json.NewEncoder(buf).Encode(q); err != nil {
		t.Fatal(err)
	}

	db = store.NewFakeStore()

	req := httptest.NewRequest("POST", "http://example.com/quotes/new", buf)
	w := httptest.NewRecorder()
	createHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 status got: %v", resp.StatusCode)
	}

	var gotQ store.Quote
	if err := json.Unmarshal(body, &gotQ); err != nil {
		t.Errorf("expected quote response got: %v", string(body))
	}
}
