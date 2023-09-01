package main

import (
	//"fmt"

	//"github.com/go-chi/chi"
	_ "github.com/lib/pq" // PostgreSQL driver
)


type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Cost   float64 `json:"cost"`
	DownloadURL string  `json:"download_url"`
}
type CartItem struct {
    BookID   int     `json:"bookid"`
    BookName string  `json:"book_name"`
    Author   string  `json:"author"`
    Genre    string  `json:"genre"`
    Cost     float64 `json:"cost"`
}

type BuyBooksResponse struct {
    Code          int
    Recommendation string
}

type User struct {
	ID       int    // Assuming your user has an ID field
	Username string // Assuming your user has a Username field
	Email    string // Assuming your user has an Email field
	// Add other fields you need
}