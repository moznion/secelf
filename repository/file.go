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

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (file_name) VALUES (?)", fileTableName))
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(fileName)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

type SearchingResult struct {
	ID       int64
	FileName string
}

func (repo *FileRepository) Search(q string) ([]*SearchingResult, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return []*SearchingResult{&SearchingResult{
			ID:       -1,
			FileName: "",
		}}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, file_name FROM %s WHERE file_name LIKE ? ORDER BY id DESC", fileTableName))
	if err != nil {
		return []*SearchingResult{&SearchingResult{
			ID:       -1,
			FileName: "",
		}}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + q + "%")
	if err != nil {
		return []*SearchingResult{&SearchingResult{
			ID:       -1,
			FileName: "",
		}}, err
	}
	defer rows.Close()

	results := make([]*SearchingResult, 0)
	for rows.Next() {
		var id int64
		var fileName string
		err = rows.Scan(&id, &fileName)
		if err != nil {
			continue
		}
		results = append(results, &SearchingResult{
			ID:       id,
			FileName: fileName,
		})
	}

	return results, nil
}
