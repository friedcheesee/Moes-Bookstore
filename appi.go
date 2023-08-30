package main

import (
	"fmt"
	//"log"
	"net/http"
	//"time"
    //"strconv"
	//"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

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