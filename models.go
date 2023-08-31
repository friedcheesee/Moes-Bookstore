package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"fmt"
	"net/http"

	//"github.com/go-chi/chi"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func searchBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("query")   // Search query
	genre := r.URL.Query().Get("genre")   // Genre filter
	author := r.URL.Query().Get("author") // Author filter
	fmt.Println(query, genre, author)

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

func searchBooks(db *sql.DB, query, genre, author string) ([]Book, error) {
	// Implement your searchBooks function logic here
	// Use the provided query, genre, and author parameters to query the database

	// Example code (modify as per your database schema and queries)
	rows, err := db.Query("SELECT bookid, book_name, author, genre, cost FROM books WHERE book_name ILIKE $1 AND genre ILIKE $2 AND author ILIKE $3",
		"%"+query+"%", "%"+genre+"%", "%"+author+"%")
	if err != nil {
		CheckError(err)
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Genre, &book.Cost)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Cost   float64 `json:"cost"`
}
