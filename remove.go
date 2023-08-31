package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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
func reactivate(db *sql.DB,email string,password string){
	authenticateUser(db, email, password)
	_, err := db.Exec("update users set active=true where email=$1", email)
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


func isAccountExpired(db *sql.DB, uid int) (bool, error) {
    var expiryDate time.Time

    err := db.QueryRow("SELECT expiry FROM users WHERE uid = $1", uid).Scan(&expiryDate)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, fmt.Errorf("User with UID %d not found", uid)
        }
        log.Println("Error querying user's expiry:", err)
        return false, err
    }

    currentDate := time.Now()
    return currentDate.After(expiryDate), nil
}
