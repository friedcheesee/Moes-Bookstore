package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "joe"
)

func main() {
	//conecct
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	//fmt.Println(db)
	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
	rows, err := getdata(db)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&name)
		CheckError(err)
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
