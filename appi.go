package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"time"
	"strconv"
	//"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if UID == 0 {
			http.Error(w, `{"status": "error", "message": "Please log in to access this endpoint"}`, http.StatusUnauthorized)
			return
		}
		// User is authenticated, call the next handler
		next.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (username and password)
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Call the login function

	isLoggedIn, err, code := logindb(db, email, password)
	if err != nil {
		// Set custom HTTP status and error message based on the error code
		httpStatus, errorMessage := getErrorDetails(code)
		http.Error(w, `{"status": "error", "message": "`+errorMessage+`"}`, httpStatus)
		return
	}
	UID = getuserid(db, email)
	// Return the login status as JSON response
	if isLoggedIn {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Login successful")
		fmt.Fprintf(w, `{"status": "success", "message": "Login successful"}`)
	} else {
		// Adjust this section to provide specific error messages based on the code
		errorMessage := "Login failed"
		if code == 1 {
			errorMessage = "User not found"
		} else if code == 2 {
			errorMessage = "Wrong password"
		}
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"status": "error", "message": "`+errorMessage+`"}`)
	}
}

// Utility function to get custom HTTP status and error message based on code
func getErrorDetails(code int) (int, string) {
	switch code {
	case 1:
		return http.StatusNotFound, "Email not found"
	case 2:
		return http.StatusUnauthorized, "Wrong password"
	case 3:
		return http.StatusUnauthorized, "Authentication failed"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	code, err := reguser(db, email, password, username)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	if code == 1 {
		http.Error(w, `{"status": "error", "message": "User already exists"}`, http.StatusBadRequest)
		return
	}
	if code == 2 {
		http.Error(w, `{"status": "error", "message": "Internal error"}`, http.StatusBadRequest)
		return
	} else {
		fmt.Fprintf(w, `{"status": "success", "message": "User registered successfully: %s"}`, username)
	}
}
func searchBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("query")   // Search query
	genre := r.URL.Query().Get("genre")   // Genre filter
	author := r.URL.Query().Get("author") // Author filter

	// Create a database connection
	// Call your searchBooks function with the provided parameters
	books, err := searchBooks(db, query, genre, author)
	if err != nil {
		http.Error(w, "Search error", http.StatusInternalServerError)
		log.Println("Error searching for books:", err)
		return
	}

	// Marshal books to JSON
	booksJSON, err := json.Marshal(books)
	if err != nil {
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}

	// Set the response content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(booksJSON)
}

func addToCartHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (bookID)
	bookIDStr := r.FormValue("bookid")
	bookid, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid bookID"}`, http.StatusBadRequest)
		return
	}
	// Check if the user is authenticated
	if UID == 0 {
		http.Error(w, `{"status": "error", "message": "Please log in to add items to your cart"}`, http.StatusUnauthorized)
		return
	}
	// Call the addToCart function with the authenticated user's UID and the bookID
	code, err := addToCart(db, UID, bookid)
    //fmt.Println("code",code)
	if err != nil {
			http.Error(w, `{"status": "error", "message": "Internal error"}`, http.StatusInternalServerError)
			return
		}
        if code == 1 {
            http.Error(w, `{"status": "error", "message": "Book already in cart"}`, http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status": "success", "message": "bookid %s added to cart successfully"}`, bookIDStr)
	}
	// Return success response

    func viewOwnedBooksHandler(w http.ResponseWriter, r *http.Request) {
        // Check if the user is authenticated
        if UID == 0 {
            http.Error(w, `{"status": "error", "message": "Please log in to view your owned books"}`, http.StatusUnauthorized)
            return
        }
        // Call the function to get the owned books for the authenticated user
        ownedBooks, err := viewOwnedBooks(db, UID)
        if err != nil {
            http.Error(w, `{"status": "error", "message": "Failed to retrieve owned books"}`, http.StatusInternalServerError)
            return
        }
        // Return the owned books as a JSON response
        responseJSON, err := json.Marshal(ownedBooks)
        if err != nil {
            http.Error(w, `{"status": "error", "message": "Failed to format response"}`, http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(responseJSON)
    }
    