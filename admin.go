package main

import (
	"database/sql"
	"fmt"
	//"fmt"
	"log"

	_ "github.com/lib/pq"
)

func removeBook(db *sql.DB, bookID int) error {
	_, err := db.Exec("DELETE FROM books WHERE bookid = $1", bookID)
	if err != nil {
		log.Println("Error removing book:", err)
		return err
	}
	log.Println("Book removed successfully")
	return nil
}

func addBook(db *sql.DB, bookName, author, genre string, cost float64) error {
	fmt.Println("Adding book")
	fmt.Println(bookName, author, genre, cost)
	_, err := db.Exec("INSERT INTO books (book_name, author, genre, cost) VALUES ($1, $2, $3, $4)",
		bookName, author, genre, cost)
	if err != nil {
		log.Println("Error adding book:", err)
		return err
	}
	log.Println("Book added successfully")
	return nil
}
func getUsers(db *sql.DB) ([]User, error) {
	// Query the database to select all active users
	rows, err := db.Query("SELECT uid, username, email FROM users WHERE active = true")
	if err != nil {
		log.Println("Error retrieving active users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func displayAvailableBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT bookid, book_name, author, genre, cost FROM books")
	if err != nil {
		log.Println("Error retrieving available books:", err)
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var bookID int
		var bookName, author, genre string
		var cost float64
		if err := rows.Scan(&bookID, &bookName, &author, &genre, &cost); err != nil {
			log.Println("Error retrieving available books:", err)
			return nil, err
		}
		book := Book{
			ID:     bookID,
			Name:   bookName,
			Author: author,
			Genre:  genre,
			Cost:   cost,
		}

		books = append(books, book)
	}
	return books, nil
}

func checkUserAdminStatus(db *sql.DB, uid int) (bool, error) {
    var isAdmin bool
    err := db.QueryRow("SELECT admin FROM users WHERE uid = $1", uid).Scan(&isAdmin)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil // User not found, assume not an admin
        }
        return false, err
    }
    return isAdmin, nil
}