package repository

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "test.sqlite3"

func TestPutShouldSuccessfully(t *testing.T) {
	scheme, _ := ioutil.ReadFile("../sql/000-file.sql")
	db, _ := sql.Open("sqlite3", dbPath)
	defer os.Remove(dbPath)
	db.Exec(string(scheme))
	db.Close()

	repo := NewFileRepository(dbPath)

	id, err := repo.Put("test")
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if id <= 0 {
		t.Errorf("cannot get id certainly")
	}
}

func TestSearchSuccessfully(t *testing.T) {
	scheme, _ := ioutil.ReadFile("../sql/000-file.sql")
	db, _ := sql.Open("sqlite3", dbPath)
	defer os.Remove(dbPath)
	db.Exec(string(scheme))
	db.Close()

	repo := NewFileRepository(dbPath)

	repo.Put("foo")
	repo.Put("bar")
	repo.Put("foobar")

	results, err := repo.Search("foo")
	if len(results) < 2 {
		t.Errorf("result is missing")
	}
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if results[0].FileName != "foobar" || results[1].FileName != "foo" {
		t.Errorf("result and/or order is wrong")
	}

	results, err = repo.Search("bar")
	if len(results) < 2 {
		t.Errorf("result is missing")
	}
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if results[0].FileName != "foobar" || results[1].FileName != "bar" {
		t.Errorf("result and/or order is wrong")
	}
}
