package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func adminconnect() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	//defer db.Close()
	err = db.Ping()
	CheckError(err)
	fmt.Println("Connected!")
	return db
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func logindb(db *sql.DB, username string, password string) {
	err := authenticateUser(db, username, password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	} else {
		fmt.Println("Authentication successful")
	}
}

func authenticateUser(db *sql.DB, username, password string) error {
	var storedPasswordHash string
	row := db.QueryRow("SELECT password FROM users WHERE username = $1", username)
	CheckError(row.Scan(&storedPasswordHash))
	// Compare stored password hash with provided hashed password
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	CheckError(err)
	return nil
}

func reguser(db *sql.DB, username string, password string) {
	hashedPassword, err := hashPassword(password)
	if userExists(db, username) {
		fmt.Println("User already exists.")
		return
	}
	CheckError(err)
	err = storeCredentials(db, username, hashedPassword)
	CheckError(err)
}

func storeCredentials(db *sql.DB, username, hashedPassword string) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

func userExists(db *sql.DB, username string) bool {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username)
	CheckError(row.Scan(&count))
	return count > 0
}
