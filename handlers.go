package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	moelog "moe/log"
	ah "moe/login"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	//"github.com/gorilla/sessions"
)

// Middleware to check if the user is an admin
func CheckUserAdminStatus(db *sql.DB, uid int) (bool, error) {
	var isAdmin bool
	err := db.QueryRow("SELECT admin FROM users WHERE uid = $1", uid).Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found, assume not an admin
		}
		return false, err
	}
	log.Println("User trying to login is admin:", isAdmin)
	return isAdmin, nil
}

// middleware handler to check if the user is logged in
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user's UID from the session cookie
		session, _ := store.Get(r, "session-name")
		uid, ok := session.Values["uid"].(int)
		//check if the user is logged in
		if !ok {
			http.Error(w, `{"status": "error", "message": "Please log in to access this endpoint"}`, http.StatusUnauthorized)
			return
		}

		// using context to pass the uid to the handler
		ctx := context.WithValue(r.Context(), "uid", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handler to aid in login
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//Parsing input parameters from the request (username and password)
	email := r.FormValue("email")
	password := r.FormValue("password")
	UID := ah.GetID(db, email) // getting uid from db (only once, to initiate session)

	//Create a session and set the user's UID with the session
	session, _ := store.Get(r, "session-name")
	session.Values["uid"] = UID
	session.Save(r, w)
	//Now that session is set, we can get uid in any other function/handler from this variable

	isLoggedIn, err, code := ah.Logindb(db, email, password)
	if err != nil {
		// Set custom HTTP status and error message based on the error code
		httpStatus, errorMessage := GetErrorDetails(code)
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
	moelog.LogEvent("logged in - handler")
}

// Utility function to get custom HTTP status and error message based on code - used for login handler
func GetErrorDetails(code int) (int, string) {
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

// to register a user, returning appropriate status codes
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
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
	//getting variables to pass to reguser function as parameters
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// Call the reguser function to register the user
	//codes: 1 for user already exists, 2 for internal error
	code, err := ah.Reguser(db, email, password, username)
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
	moelog.LogEvent("registered user - handler")
}

// to search for books using query, genre and author
func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters
	query := r.URL.Query().Get("query")
	genre := r.URL.Query().Get("genre")
	author := r.URL.Query().Get("author")

	// Create a database connection
	// Call your searchBooks function with the provided parameters
	books, err := searchBooks(db, query, genre, author)
	if err != nil {
		http.Error(w, "Search error", http.StatusInternalServerError)
		log.Println("Error searching for books:", err)
		return
	}

	// using marshal to convert books struct to json (defined in models.go)
	booksJSON, err := json.Marshal(books)
	if err != nil {
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}

	// Set the response content type and write JSON response
	moelog.LogEvent("searched books - handler")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(booksJSON)
}

// to add to cart, returning appropriate status codes
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (bookID)
	bookIDStr := r.FormValue("bookid")
	bookid, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid bookID"}`, http.StatusBadRequest)
		return
	}
	//getting uid from session
	UID := r.Context().Value("uid").(int)
	// Check if the user is authenticated
	if UID == 0 {
		http.Error(w, `{"status": "error", "message": "Please log in to add items to your cart"}`, http.StatusUnauthorized)
		return
	}

	// Call the addToCart function with the authenticated user's UID and the bookID
	code, err := AddToCart(db, UID, bookid)
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
	moelog.LogEvent("added to cart - handler")
}

// to view owned books
func ViewOwnedBooksHandler(w http.ResponseWriter, r *http.Request) {
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
	// Check if no books are owned
	if len(ownedBooks) == 0 {
		// No books are owned, return a "success" message
		response := map[string]interface{}{
			"status":  "success",
			"message": "No books owned",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	// Return the owned books as a JSON response
	responseJSON, err := json.Marshal(ownedBooks)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to format response"}`, http.StatusInternalServerError)
		return
	}
	moelog.LogEvent("viewed owned books - handler")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// to view reviews, and return appropriate status codes
// codes 1 for review already exists, 2 for internal error
func GiveReviewHandler(w http.ResponseWriter, r *http.Request) {
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
	moelog.LogEvent("gave review - handler")
}

// handler to delete bookid from cart
func DeleteFromCartHandler(w http.ResponseWriter, r *http.Request) {

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
	fmt.Fprintf(w, `{"status": "success", "message": "Book deleted from cart successfully. if present"}`)
	moelog.LogEvent("deleted from cart - handler")
}

// handler to buy books
// codes 1 for conflicting books, 2 for internal error
func BuyBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
	UID := r.Context().Value("uid").(int)
	if UID == 0 {
		http.Error(w, `{"status": "error", "message": "Please log in to buy books"}`, http.StatusUnauthorized)
		return
	}

	// Call the function to buy books
	code, err, recommendations := buyBooks(db, UID)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to buy books"}`, http.StatusInternalServerError)
		return
	}

	if code == 1 {
		http.Error(w, `{"status": "error", "message": "Remove conflicting book from the cart to buy books"}`, http.StatusBadRequest)
		return
	}

	// Check if there are recommendations
	if len(recommendations) > 0 {
		recommendation := recommendations[0]

		// Construct the response
		response := fmt.Sprintf(`{"status": "success", "message": "Books bought successfully, if present in cart", "recommendation": {"name": "%s", "author": "%s"}}`, recommendation.Name, recommendation.Author)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	} else {
		// No recommendations available
		response := `{"status": "success", "message": "Books bought successfully, if present in cart", "recommendation": null}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		moelog.LogEvent("bought books - handler")
	}
}

// middle ware handler to view cart
func ViewCartHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user's UID from the global variable or session
	uid := r.Context().Value("uid").(int)

	// Call the function to retrieve items from the cart
	cartItems, err := viewCart(db, uid)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Error fetching cart items"}`, http.StatusInternalServerError)
		return
	}

	// Check if the cart is empty
	if len(cartItems) == 0 {
		// Cart is empty, return a status message
		response := map[string]interface{}{
			"status":  "success",
			"message": "Cart is empty",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return the cart items as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
	moelog.LogEvent("viewed cart - handler")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong!")
	moelog.LogEvent("got pinged - handler")
}

// middleware handler to delete account
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (email and password)
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Call the deactivate function
	ah.Delete(db, email, password)
	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "success", "message": "Account deactivated successfully"}`)
	moelog.LogEvent("deleted account - handler")
}

// middleware handler to logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["uid"] = nil //deletes session linked to uid
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "message": "Logout successful"}`)

	moelog.LogEvent("logged out - handler")
}

// middleware handler to check if account is active(not deleted)
func CheckActiveAccount(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Authenticate the user
		isAuthenticated, err := ah.AuthenticateUser(db, email, password)
		if err != nil || !isAuthenticated {
			http.Error(w, `{"status": "error", "message": "Invalid credentials/User doesn't exist"}`, http.StatusUnauthorized)
			return
		}

		// Check if the account is active
		isActive := ah.IsAccountActive(db, email)
		if !isActive {
			http.Error(w, `{"status": "error", "message": "Account has been deleted, please sign up with another email"}`, http.StatusUnauthorized)
			return
		}
		// User is authenticated and account is active, call the next handler
		next.ServeHTTP(w, r)
	}
}

// to check if the user is admin
func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user's UID from the session or wherever it's stored
		session, _ := store.Get(r, "session-name")
		uid, ok := session.Values["uid"].(int)
		if !ok {
			http.Error(w, `{"status": "error", "message": "Please log in to access this endpoint"}`, http.StatusUnauthorized)
			return
		}

		// Query the database to check if the user with the given UID is an admin
		isAdmin, err := CheckUserAdminStatus(db, uid)
		if err != nil {
			http.Error(w, `{"status": "error", "message": "Error checking admin status"}`, http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.Error(w, `{"status": "error", "message": "You do not have admin privileges"}`, http.StatusForbidden)
			return
		}

		// If the user is an admin, call the next handler.
		next.ServeHTTP(w, r)
	})
}

// to add books for admins
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (bookName, author, genre, cost, stock)
	bookName := r.FormValue("bookName")
	author := r.FormValue("author")
	genre := r.FormValue("genre")
	costStr := r.FormValue("cost")
	// Convert cost and stock to float64 and int respectively
	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid cost"}`, http.StatusBadRequest)
		return
	}

	// Call the addBook function
	err = ah.AddBook(db, bookName, author, genre, cost)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to add book, book already exists"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "message": "Book added successfully"}`)
	moelog.LogEvent("added book - handler")
}

// to remove books for admins
func RemoveBookHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input parameters from the request (bookID)
	bookIDStr := r.FormValue("bookID")

	// Convert bookID to integer
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid book ID"}`, http.StatusBadRequest)
		return
	}

	// Call the removeBook function
	err = ah.RemoveBook(db, bookID)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to remove book"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "message": "Book removed successfully, if present"}`)
	moelog.LogEvent("removed book - handler")
}

// to view users for admins
func ViewUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Call the getUsers function
	users, err := ah.GetUsers(db)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to retrieve users"}`, http.StatusInternalServerError)
		return
	}

	// Return the users as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	moelog.LogEvent("viewed users - handler")
}

// to view all available books
func ViewAvailableBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Call the displayAvailableBooks function
	books, err := ah.DisplayAvailableBooks(db)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Failed to retrieve books"}`, http.StatusInternalServerError)
		return
	}

	// Return the books as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	moelog.LogEvent("viewed available books - handler")
}
