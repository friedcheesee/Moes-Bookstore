package main

import (
	//"fmt"
	"database/sql"
	"log"
	"net/http"
	"os"
    "github.com/gorilla/sessions"
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
var store = sessions.NewCookieStore([]byte("12345678"))

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
	r.Post("/login", checkActiveAccount(loginHandler))
	r.Post("/reguser", registerUserHandler)
	r.Route("/user", func(r chi.Router){
		r.Use(authenticate)
		r.Get("/search", searchBooksHandler)
		r.Post("/cart/add", addToCartHandler)
		r.Post("/cart/delete", deleteFromCartHandler)
		r.Post("/cart/buy", buyBooksHandler)
		r.Post("/cart/view", viewCartHandler)
		r.Post("/inventory", viewOwnedBooksHandler)
		r.Post("/review", giveReviewHandler)
		r.Post("/delete", deactivateHandler)
		r.Post("/logout", logoutHandler)
	})
	
	r.Route("/admin", func(r chi.Router){
		r.Use(isAdmin)
		r.Post("/add", addBookHandler)
		r.Post("/delete", removeBookHandler)
		r.Post("/view", viewUsersHandler)
		r.Post("/view/books", viewAvailableBooksHandler)
	})

	//r.Post("/available", authenticate(displayAvailableBooksHandler))
	//r.Post("/reviews", authenticate(displayBookReviewsHandler))
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
	////reguser(db, email, password, username)-
	////logindb(db, email, password)-
	//displayBookReviews(db, 1)
	////bookname:="Sample"
	////searchBooks(db,bookname, "", "")-
	////ff := 
	////fmt.Println(ff)
	//deactivate(db, username, password)
	//addToCart(db, ff, 1)-
	//viewOwnedBooks(db, ff)-
	//giveReview(db, ff, 1, review)--
	//displayAvailableBooks(db)
	//deleteFromCart(db, ff, 1)--
	//err := buyBooks(db, ff)-
	//CheckError(err)
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
