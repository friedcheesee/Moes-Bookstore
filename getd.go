package main
import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // postgres golang driver
)
func getdata(db *sql.DB) (*sql.Rows, error) {
	query := "SELECT name,roll,uid FROM joe"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func printnames(rows *sql.Rows) {
	defer rows.Close()
	for rows.Next() {
		var uid,roll int
		var name string
		err := rows.Scan(&name,&roll,&uid)
		CheckError(err)
		fmt.Printf("Name: %s\nRoll: %d\nID: %d\n", name,roll,uid)
	}
}