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

type FileRow struct {
	ID       int64  `json:"id"`
	FileName string `json:"file_name"`
}

func (repo *FileRepository) Single(id int64) (*FileRow, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("SELECT file_name FROM %s WHERE id = ?", fileTableName))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var fileName string
	err = stmt.QueryRow(id).Scan(&fileName)
	if err != nil {
		return nil, err
	}

	return &FileRow{
		ID:       id,
		FileName: fileName,
	}, nil
}

func (repo *FileRepository) Search(q string) ([]*FileRow, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return []*FileRow{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, file_name FROM %s WHERE file_name LIKE ? ORDER BY id DESC", fileTableName))
	if err != nil {
		return []*FileRow{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + q + "%")
	if err != nil {
		return []*FileRow{}, err
	}
	defer rows.Close()

	results := make([]*FileRow, 0)
	for rows.Next() {
		var id int64
		var fileName string
		err = rows.Scan(&id, &fileName)
		if err != nil {
			continue
		}
		results = append(results, &FileRow{
			ID:       id,
			FileName: fileName,
		})
	}

	return results, nil
}
