package main

import (
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq"
)

func delconnect() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=deletor password=det dbname=%s sslmode=disable", host, port, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	err = db.Ping()
	CheckError(err)
	fmt.Println("Connected!")
	return db
}

func deactivate(db *sql.DB,email string,password string){
	authenticateUser(db, email, password)
	_, err := db.Exec("update users set active=false where email=$1", email)
	CheckError(err)
	fmt.Println("User account deactivated")
}
