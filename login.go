package main

import (
	"database/sql"
	"fmt"

	//"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

func adminconnect() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database //dont run this if its in a seaparte function...
	//defer db.Close()

	//fmt.Println(db)
	// check db
	err = db.Ping()
	CheckError(err)
	fmt.Println("Connected!")
	return db
}
func printnames(rows *sql.Rows) {

	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&name)
		CheckError(err)
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
}
