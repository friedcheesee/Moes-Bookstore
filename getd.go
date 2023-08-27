package main

import (
	"database/sql"

	_ "github.com/lib/pq" // postgres golang driver
)

func getdata(db *sql.DB) (*sql.Rows, error) {
	query := "SELECT name FROM joe"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// func getdata(db *sql.DB) (*sql.Rows, error) {
// 	rows, err := db.Query(`SELECT name, Roll FROM joe`)
// 	CheckError(err)

// 	defer rows.Close()
// 	for rows.Next() {
// 		var name string
// 		var roll int
// 		err = rows.Scan(&name, &roll)
// 		CheckError(err)

// 		fmt.Println(name, roll)
// 	}

// 	CheckError(err)
// }
