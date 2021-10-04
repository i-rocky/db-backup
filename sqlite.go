package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Record struct {
	Id int
	Name string
	File string
	Time time.Time
}

type SQLite struct {
	db *sql.DB
}

func (s *SQLite) connect() error {
	db, err := sql.Open("sqlite3", "file:local.db")
	if err != nil {
		return err
	}

	s.db = db
	s.db.SetMaxOpenConns(1)
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS records (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, file TEXT, time TEXT)")
	return err
}

func (s *SQLite) record(file string) error {
	_, err := s.db.Exec("INSERT INTO records (name, file, time) VALUES (?, ?, ?)", time.Now().Format(time.RFC1123), file, time.Now().Format(time.RFC1123Z))
	return err
}

func (s *SQLite) records() ([]Record, error) {
	var list []Record
	rows, err := s.db.Query("SELECT * FROM records")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record Record
	for rows.Next() {
		var t string
		err := rows.Scan(&record.Id, &record.Name, &record.File, &t)
		if err != nil {
			continue
		}
		tm, err := time.Parse(time.RFC1123Z, t)
		if err != nil {
			continue
		}
		record.Time = tm
		list = append(list, record)
	}

	return list, err
}

func (s *SQLite) remove(ids ...int) error {
	_, err := s.db.Exec("DELETE FROM records WHERE id IN (?)", ids)
	return err
}
