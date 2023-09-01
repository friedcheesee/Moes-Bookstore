package ah
import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	_ "github.com/lib/pq"        // PostgreSQL driver
	"golang.org/x/crypto/bcrypt" // For encrypting and decrypting passwords
	"moe/log"
)



//global database connection
func Adminconnect() *sql.DB {

	// Access environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	// open database
	db, err := sql.Open("postgres", psqlconn)
	moelog.CheckError(err)
	
	// check connection
	err = db.Ping()
	moelog.CheckError(err)
	fmt.Println("Connected!")
	moelog.LogEvent("Connected to database")
	return db
}

//function will will convert the password into a hash, which will be stored in the database.
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

//to fetch ID from database, used only once during login function to set cookies.
func GetID(db *sql.DB, email string) int {
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

// Codes returned by this function, to debug/show status of login
// 0 - success
// 1 - email not registered
// 2 - wrong password
func Logindb(db *sql.DB, email string, password string) (bool, error, int) {
	if !userExists(db, email) {
		moelog.LogEvent("User not found")
		return false, nil, 1 // User not found
	}
	isAuthenticated, err := AuthenticateUser(db, email, password)
	if err != nil {
		moelog.CheckError(err)
		moelog.LogEvent("Authentication error")
		return false, err, 2 // Authentication error
	}
	if !isAuthenticated {
		moelog.LogEvent("Incorrect password")
		return false, nil, 3 // Incorrect password
	}
	moelog.LogEvent("User logged in successfully")
	return true, nil, 0 // Success
}

// Authenticate user using the provided email and password
func AuthenticateUser(db *sql.DB, email, password string) (bool, error) {
	var storedPasswordHash string
	row := db.QueryRow("SELECT password FROM users WHERE email = $1", email)
	moelog.CheckError(row.Scan(&storedPasswordHash))

	// Compare stored password hash with provided hashed password
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	if err != nil {
		log.Println("Wrong password", err)
		return false, err
	}
	return true, nil
}

// codes returned by this function, to debug/show status of registration
// 0 success
// 1 user exists
// 2 internal error
func Reguser(db *sql.DB, email, password, username string) (int, error) {
	hashedPassword, err := hashPassword(password)
	if userExists(db, email) {
		fmt.Println("User already exists.")
		moelog.LogEvent("User already exists")
		return 1, nil
	}
	if err != nil {
		moelog.CheckError(err)
		moelog.LogEvent("Error hashing password")
		fmt.Println("Error hashing password")
		return 2, err
	}
	if !validateEmail(email) {
		moelog.LogEvent("Invalid email")
		fmt.Println("Invalid email")
		return 2, nil
	}
	storeCredentials(db, username, hashedPassword, email)
	moelog.CheckError(err)
	return 0, nil
}
//to store credentials in db while registering user
func storeCredentials(db *sql.DB, username string, hashedPassword string, email string) {
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, hashedPassword, email)
	moelog.CheckError(err)
	moelog.LogEvent("User registered successfully")
	fmt.Println("User registered successfully")
}

//used in reguser and loginuser to check if user exists
func userExists(db *sql.DB, email string) bool {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email)
	moelog.CheckError(row.Scan(&count))
	return count > 0
}

//used in reguser to check if email is valid
func validateEmail(email string) bool {
	// Regular expression pattern for basic email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Use the regular expression to match against the email
	return re.MatchString(email)
}

//used in handlers to check if user logging in has deleted their account
func IsAccountActive(db *sql.DB, email string) bool {
	var active bool
	err := db.QueryRow("SELECT active FROM users WHERE email = $1 ", email).Scan(&active)
	if err != nil {
		return false
	}
	return active
}

func Delete(db *sql.DB,email string,password string){
	AuthenticateUser(db, email, password)
	_, err := db.Exec("update users set active=false where email=$1", email)
	moelog.CheckError(err)
	fmt.Println("User account deactivated")
}
