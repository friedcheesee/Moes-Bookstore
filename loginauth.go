package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

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

func reguser(db *sql.DB, username string, password string,email string) {
	hashedPassword, err := hashPassword(password)
	if userExists(db, email) {
		fmt.Println("User already exists.")
		return
	}
	CheckError(err)
	if(!validateEmail(email)){
		fmt.Println("Invalid email")
		log.Fatal("Invalid email")}
	storeCredentials(db, username, hashedPassword,email)
	CheckError(err)
}

func storeCredentials(db *sql.DB, username string, hashedPassword string,email string){
	var maxUID int
    err := db.QueryRow("SELECT COALESCE(MAX(uid),0) FROM users").Scan(&maxUID)
    CheckError(err)
    uid := maxUID + 1
	_,err = db.Exec("INSERT INTO users (username, password, uid, active,email) VALUES ($1, $2, $3, TRUE,$4)", username, hashedPassword, uid,email)
	CheckError(err)
	fmt.Println("User registered successfully")
}

func userExists(db *sql.DB, email string) bool {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email)
	CheckError(row.Scan(&count))
	return count > 0
}
func validateEmail(email string) bool {
	// Regular expression pattern for basic email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	
	// Compile the regular expression
	re := regexp.MustCompile(pattern)
	
	// Use the regular expression to match against the email
	return re.MatchString(email)
}

func isUserActive(db *sql.DB, uid int) (bool) {
	var isActive bool
	err := db.QueryRow("SELECT active FROM users WHERE uid = $1", uid).Scan(&isActive)
	CheckError(err)
	return isActive
}