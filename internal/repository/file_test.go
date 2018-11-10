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
	scheme, _ := ioutil.ReadFile("../../sql/000-file.sql")
	db, _ := sql.Open("sqlite3", dbPath)
	defer os.Remove(dbPath)
	db.Exec(string(scheme))
	db.Close()

	repo := NewFileRepository(dbPath)

	id, err := repo.Put("test", "salt")
	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}
	if id <= 0 {
		t.Errorf("cannot get id certainly")
	}
}

func TestSearchSuccessfully(t *testing.T) {
	scheme, _ := ioutil.ReadFile("../../sql/000-file.sql")
	db, _ := sql.Open("sqlite3", dbPath)
	defer os.Remove(dbPath)
	db.Exec(string(scheme))
	db.Close()

	repo := NewFileRepository(dbPath)

	repo.Put("foo", "salt1")
	repo.Put("bar", "salt2")
	repo.Put("foobar", "salt3")

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

	results, err = repo.Search("qux")
	if len(results) != 0 {
		t.Errorf("result is wrong")
	}
	if err != nil {
		t.Errorf("got unexpected error")
	}
}

func TestSingleSuccessfully(t *testing.T) {
	scheme, _ := ioutil.ReadFile("../../sql/000-file.sql")
	db, _ := sql.Open("sqlite3", dbPath)
	defer os.Remove(dbPath)
	db.Exec(string(scheme))
	db.Close()

	repo := NewFileRepository(dbPath)

	id1, _ := repo.Put("foo", "salt1")
	id2, _ := repo.Put("bar", "salt2")
	id3, _ := repo.Put("foobar", "salt3")

	row, err := repo.Single(id1)
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if row.FileName != "foo" {
		t.Errorf("unexpected result [got:%s, expected:%s]", row.FileName, "foo")
	}

	row, err = repo.Single(id2)
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if row.FileName != "bar" {
		t.Errorf("unexpected result [got:%s, expected:%s]", row.FileName, "bar")
	}

	row, err = repo.Single(id3)
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if row.FileName != "foobar" {
		t.Errorf("unexpected result [got:%s, expected:%s]", row.FileName, "foobar")
	}

	row, err = repo.Single(0)
	if err == nil {
		t.Errorf("missing expected error")
	}
	if row != nil {
		t.Errorf("result is not nil")
	}
}
