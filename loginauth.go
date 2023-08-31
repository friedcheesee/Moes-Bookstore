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
func getID(db *sql.DB, email string) (int) {
    var uid int
    err := db.QueryRow("SELECT uid FROM users WHERE email = $1", email).Scan(&uid)
    if err != nil {
        if err == sql.ErrNoRows {
			log.Printf("No user found with the provided email")
            return 0 
        }
        return 0
    }
    return uid
}


// 0 - success
// 1 - email not registered
// 2 - wrong password
func logindb(db *sql.DB, email string, password string) (bool, error, int) {
	if !userExists(db, email) {
		fmt.Println("User not found")
		return false, nil, 1 // User not found
	}
	isAuthenticated, err := authenticateUser(db, email, password)
	if err != nil {
		fmt.Println("wring pw",err)
		return false, err, 2 // Authentication error
	}
	if !isAuthenticated {
		fmt.Println("wrong pw")
		return false, nil, 3 // Incorrect password
	}
	return true, nil, 0 // Success
}


func authenticateUser(db *sql.DB, email, password string) (bool,error) {
	var storedPasswordHash string
	row := db.QueryRow("SELECT password FROM users WHERE email = $1", email)
	CheckError(row.Scan(&storedPasswordHash))
	// Compare stored password hash with provided hashed password
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	if err != nil {
		log.Println("Wrong password",err)
		return false, err
	}
	return true,nil
}
// 0 success
// 1 user exists
// 2 internal error
func reguser(db *sql.DB, email , password ,username string)(int ,error) {
	hashedPassword, err := hashPassword(password)
	if userExists(db, email) {
		fmt.Println("User already exists.")
		return 1,nil
	}
	if err != nil {
		log.Print("error registering user",err)
		return 2,err
	}
	if !validateEmail(email) {
		fmt.Println("Invalid email")
		log.Fatal("Invalid email")
		return 2,nil
	}
	storeCredentials(db, username, hashedPassword, email)
	CheckError(err)
	return 0,nil
}


func storeCredentials(db *sql.DB, username string, hashedPassword string, email string) {
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, hashedPassword, email)
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


func isAccountActive(db *sql.DB, email, password string) bool {
    var active bool
    err := db.QueryRow("SELECT active FROM users WHERE email = $1 AND password = $2", email, password).Scan(&active)
    if err != nil {
        return false
    }
    return active
}
