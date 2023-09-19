package ah

import (
	"database/sql"
	"fmt"
	"log"
	"moe/log"
	"moe/models"
)

// function for admins to remove a book from available books
func RemoveBook(db *sql.DB, bookID int) error {
	// query to delete book based on bookid
	_, err := db.Exec("DELETE FROM books WHERE bookid = $1", bookID)
	if err != nil {
		log.Println("Error removing book:", err)
		return err
	}
	moelog.LogEvent("(admin) Book removed successfully, if present")
	return nil
}

// function for admins to add a book to available books
func AddBook(db *sql.DB, bookName, author, genre string, cost float64) error {
	fmt.Println("Adding book")
	fmt.Println(bookName, author, genre, cost)
	_, err := db.Exec("INSERT INTO books (book_name, author, genre, cost) VALUES ($1, $2, $3, $4)",
		bookName, author, genre, cost)
	if err != nil {
		log.Println("Error adding book:", err)
		return err
	}
	moelog.LogEvent("(admin) Book added successfully")
	return nil
}

// function for admins to view all users
func GetUsers(db *sql.DB) ([]models.User, error) {

	// Query the database to select all active users (in other words, accounts which are not deleted yet)
	rows, err := db.Query("SELECT uid, username, email FROM users WHERE active = true")
	if err != nil {
		log.Println("Error retrieving active users:", err)
		return nil, err
	}
	defer rows.Close()

	// Store the results in an array of User structs
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// function for admins to view all available books
func DisplayAvailableBooks(db *sql.DB) ([]models.Book, error) {
	rows, err := db.Query("SELECT b.bookid, b.book_name, b.author, b.genre, b.cost, b.download_url, COALESCE(r.review, '') " +
		"FROM books b LEFT JOIN reviews r ON b.bookid = r.bookid")
	if err != nil {
		log.Println("Error retrieving available books:", err)
		return nil, err
	}
	defer rows.Close()

	// Store the results in an array of Book structs
	var books []models.Book
	for rows.Next() {
		var bookID int
		var bookName, author, genre, downloadURL, review string
		var cost float64
		if err := rows.Scan(&bookID, &bookName, &author, &genre, &cost, &downloadURL, &review); err != nil {
			log.Println("Error retrieving available books:", err)
			return nil, err
		}
		book := models.Book{
			ID:          bookID,
			Name:        bookName,
			Author:      author,
			Genre:       genre,
			Cost:        cost,
			DownloadURL: downloadURL,
			Review:      review,
		}
		books = append(books, book)
	}
	moelog.LogEvent("(admin) Available books retrieved successfully")
	return books, nil
}

