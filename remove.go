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

//call fn only after authentication
func getuserid(db *sql.DB,username string) (*sql.Rows) {
	rows, err := db.Query("select UID from users where username=$1", username)
	CheckError(err)
	return rows 
}

func deactivate(db *sql.DB,username string,password string){
	logindb(db, username, password)
	_, err := db.Exec("update users set active=0 where username=$1", username)
	CheckError(err)
	fmt.Println("User account deactivated")
}
func reactivate(db *sql.DB,username string,password string){
	logindb(db, username, password)
	_, err := db.Exec("update users set active=1 where username=$1)", username)
	CheckError(err)
	fmt.Println("User account reactivated")
}

func readID(rows *sql.Rows) {
    defer rows.Close()
    for rows.Next() {
        var userID int
        if err := rows.Scan(&userID); err != nil {
            panic(err)
        }
        fmt.Printf("User ID: %d\n", userID)
    }
    if err := rows.Err(); err != nil {
        panic(err)
    }
}