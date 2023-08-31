package main

import (
	//"fmt"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "moe"
)

var db *sql.DB
var UID int = 0

func main1() {
	db := adminconnect()
	defer db.Close()
	//username := "banana"
	email := "fried@mail.com"
	password := "abcd"

	//reguser(db, email, password, username)
	logindb(db, email, password)
}

func main() {
	db = adminconnect()
	defer db.Close()
	logFile := initiatelog()
	defer logFile.Close()
	r := chi.NewRouter()
	// Define the login route
	r.Post("/login", loginHandler)
	r.Post("/reguser", registerUserHandler)
	r.Get("/search", authenticate(searchBooksHandler))
	//write a call to add to cart
	r.Post("/addtocart", addToCartHandler)
	// Start the HTTP server
	http.ListenAndServe("localhost:8080", r)
	//conecct
	////db := adminconnect()
	////defer db.Close()
	//rows, err := getdata(db)
	//CheckError(err)
	//printnames(rows)
	//reguser(db)
	////username := "banana"
	////email:= "fried@mail.com"
	////password := "abcd"
	//review:= "this is a review2"
	////reguser(db, email, password, username)
	////logindb(db, email, password)
	//displayBookReviews(db, 1)
	////bookname:="Sample"
	////searchBooks(db,bookname, "", "")
	////ff := getuserid(db, email)
	////fmt.Println(ff)
	//deactivate(db, username, password)
	//addToCart(db, ff, 1)
	//viewOwnedBooks(db, ff)
	//giveReview(db, ff, 1, review)
	//displayAvailableBooks(db)
	//deleteFromCart(db, ff, 1)
	//err := buyBooks(db, ff)
	//CheckError(err)
	////if(!isUserActive(db,ff)){
	////	fmt.Println("User is not active")
	////}
}
func CheckError(err error) {
	if err != nil {
		log.Println("Error: &s", err)
		//panic(err)
	}
}
func initiatelog() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	CheckError(err)
	log.SetOutput(logFile)
	return logFile
}
