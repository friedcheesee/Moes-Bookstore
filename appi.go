package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    //"github.com/gorilla/sessions"
	//"time"
    "context"
	"strconv"
	//"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)
func authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "session-name")
        uid, ok := session.Values["uid"].(int)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Now you have the uid available, you can pass it to the handler using a context
        ctx := context.WithValue(r.Context(), "uid", uid)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func authenticate1(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // Get the user's UID from the global variable or session
    
		//if UID == 0 {
		//	http.Error(w, `{"status": "error", "message": "Please log in to access this endpoint"}`, http.StatusUnauthorized)
		//	return
		//}
		// User is authenticated, call the next handler
		next.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (username and password)
	email := r.FormValue("email")
	password := r.FormValue("password")
	UID:=getID(db, email)
	session, _ := store.Get(r, "session-name")
	session.Values["uid"] = UID
	session.Save(r, w)

	isLoggedIn, err, code := logindb(db, email, password)
	if err != nil {
		// Set custom HTTP status and error message based on the error code
		httpStatus, errorMessage := getErrorDetails(code)
		http.Error(w, `{"status": "error", "message": "`+errorMessage+`"}`, httpStatus)
		return
	}
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
    UID := r.Context().Value("uid").(int)
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
    UID := r.Context().Value("uid").(int)
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
func giveReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
    UID := r.Context().Value("uid").(int)
	if UID == 0 {
		http.Error(w, `{"status": "error", "message": "Please log in to give a review"}`, http.StatusUnauthorized)
		return
	}
	// Parse input parameters from the request (bookID and review)
	bookIDStr := r.FormValue("bookID")
	review := r.FormValue("review")
	// Convert bookID to integer
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid book ID"}`, http.StatusBadRequest)
		return
	}
	// Call the function to give a review
	code, err := giveReview(db, UID, bookID, review)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to give a review"}`, http.StatusInternalServerError)
		return
	}
	if code == 1 {
		http.Error(w, `{"status": "error", "message": "Review already exists"}`, http.StatusBadRequest)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "message": "Review added successfully"}`)
}
func deleteFromCartHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
    UID := r.Context().Value("uid").(int)
	if UID == 0 {
		http.Error(w, `{"status": "error", "message": "Please log in to delete a book from the cart"}`, http.StatusUnauthorized)
		return
	}

	// Parse input parameters from the request (bookID)
	bookIDStr := r.FormValue("bookID")

	// Convert bookID to integer
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid book ID"}`, http.StatusBadRequest)
		return
	}

	// Call the function to delete a book from the cart
	deleteFromCart(db, UID, bookID)

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "message": "Book deleted from cart successfully"}`)
}
func buyBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
    UID := r.Context().Value("uid").(int)
	if UID == 0 {
		http.Error(w, `{"status": "error", "message": "Please log in to buy books"}`, http.StatusUnauthorized)
		return
	}

	// Call the function to buy books
	code, err := buyBooks(db, UID)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to buy books"}`, http.StatusInternalServerError)
		return
	}
	if code == 1 {
		http.Error(w, `{"status": "error", "message": "Remove book from the cart to buy books"}`, http.StatusBadRequest)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "message": "Books bought successfully"}`)
}

func viewCartHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user's UID from the global variable or session
    UID := r.Context().Value("uid").(int)
	uid := UID

	// Call the function to retrieve items from the cart
	cartItems, err := viewCart(db, uid)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Error fetching cart items"}`, http.StatusInternalServerError)
		return
	}

	// Return the cart items as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
}

func deactivateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (email and password)
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Call the deactivate function
	deactivate(db, email, password)
	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "success", "message": "Account deactivated successfully"}`)
}
func reactivateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (email and password)
	email := r.FormValue("email")
	password := r.FormValue("password")
	reactivate(db, email, password)
	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "success", "message": "Account reactivated successfully"}`)
}
func authenticateActive(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.Header.Get("email")
		password := r.Header.Get("password")
		isActive := isAccountActive(db, email, password)
		if !isActive {
			http.Error(w, `{"status": "error", "message": "Account is not active. Please reactivate your account"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
