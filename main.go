package main

import (
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
	db := adminconnect()
	//rows, err := getdata(db)
	//CheckError(err)
	//printnames(rows)
	//reguser(db)
	username := "apple"
	password := "abcd"
	reguser(db, username, password)
	logindb(db, username, password)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
