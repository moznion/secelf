package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const fileTableName = "file"

type FileRepository struct {
	dbPath string
}

func NewFileRepository(dbPath string) *FileRepository {
	return &FileRepository{
		dbPath: dbPath,
	}
}

func (repo *FileRepository) Put(fileName string) (int64, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return -1, err
	}
	defer db.Close()
	res, err := db.Exec(fmt.Sprintf(`INSERT INTO %s (file_name) VALUES ("%s");`, fileTableName, fileName))
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
