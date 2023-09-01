package main

import (
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq"
)

func deactivate(db *sql.DB,email string,password string){
	authenticateUser(db, email, password)
	_, err := db.Exec("update users set active=false where email=$1", email)
	CheckError(err)
	fmt.Println("User account deactivated")
}
