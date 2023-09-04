package main

import (
	"database/sql"
	"log"
	"moe/log"
"net/http"
"moe/middleware"
	"os"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


// load variables from .env file
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// initialising single database connection to use throughout the program
var db *sql.DB
var logFile *os.File

// initialising cookie store, to allow multiple users to be connected to the database
var store *sessions.CookieStore

func main() {
	loadEnv()         //load variables from env to login to db
	initCookieStore() //initialise cookie store

	//initialising log file
	logFile = moelog.Initiatelog()
	defer logFile.Close()

	//connecting to database
	db=ah.Adminconnect()
	//db = ah.Newconnect()
	defer db.Close()

	//using chi router to handle requests
	r := chi.NewRouter()
	r.Post("/login", CheckActiveAccount(LoginHandler)) //if account isnt deleted, lets you login
	r.Post("/reguser", RegisterUserHandler)
	r.Post("/ping", pingHandler)
	//user is a 'subrouter'- every request to /user will undergo the authenticate middleware
	r.Route("/user", func(r chi.Router) {
		r.Use(Authenticate)
		r.Get("/search", SearchBooksHandler)
		r.Post("/cart/add", AddToCartHandler)
		r.Post("/cart/delete", DeleteFromCartHandler)
		r.Post("/cart/buy", BuyBooksHandler)
		r.Post("/cart/view", ViewCartHandler)
		r.Post("/inventory", ViewOwnedBooksHandler)
		r.Post("/review", GiveReviewHandler)
		r.Post("/delete", DeleteHandler)
		r.Post("/logout", LogoutHandler)
	})

	//admin is a 'subrouter'- every request to /admin will undergo the isAdmin middleware
	r.Route("/admin", func(r chi.Router) {
		r.Use(IsAdmin)
		r.Post("/add", AddBookHandler)
		r.Post("/delete", RemoveBookHandler)
		r.Post("/view", ViewUsersHandler)
		r.Post("/view/books", ViewAvailableBooksHandler)
	})
	http.ListenAndServe("0.0.0.0:8080", r)
}



// creates session store variable for cookies with a random key
func initCookieStore() *sessions.CookieStore {
	cookieKey := os.Getenv("COOKIE_KEY")
	// Create the CookieStore
	store = sessions.NewCookieStore([]byte(cookieKey)) // contains random key generated by "crypto/rand" ,"encoding/hex"
	return store
}
