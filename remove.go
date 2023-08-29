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
func getuserid(db *sql.DB,username string) int {
	rows, err := db.Query("select UID from users where username=$1", username)
	CheckError(err)
	return readID(rows)

}

func deactivate(db *sql.DB,username string,password string){
	logindb(db, username, password)
	_, err := db.Exec("update users set active=true where username=$1", username)
	CheckError(err)
	fmt.Println("User account deactivated")
}
func reactivate(db *sql.DB,username string,password string){
	logindb(db, username, password)
	_, err := db.Exec("update users set active=true where username=$1", username)
	CheckError(err)
	fmt.Println("User account reactivated")
}

func readID(rows *sql.Rows) int{
    defer rows.Close()
    for rows.Next() {
        var userID int
        if err := rows.Scan(&userID); err != nil {
            panic(err)
        }
        //fmt.Printf("User ID: %d\n", userID)
		return userID
    }
    if err := rows.Err(); err != nil {
        CheckError(err)
    }
	return 0
	
}
func readID1(rows *sql.Rows) (int) {
    defer rows.Close()
    if rows.Next() {
        var userID int
        if err := rows.Scan(&userID); err != nil {
			CheckError(err)
            return userID   // Return error if scanning fails
        }
        fmt.Printf("User ID: %d\n", userID)
        return userID// Return userID if scanning is successful
    }

    if err := rows.Err(); err != nil {
        CheckError(err) // Return error if there's an error in rows
    }

    	return 0  // Return specific error for no rows found
}