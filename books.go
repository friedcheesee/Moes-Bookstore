package main

import (
	"database/sql"
	"fmt"
	"log"
)

// 0 added
// 1 already exists
// 2 internal error
func addToCart(db *sql.DB, uid, bookid int) (int, error) {
	// Check if the book already exists in the cart
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid).Scan(&count)
	//check the count of statements where the uid owns bookid. if owns 1 time, cant own again.
	if err != nil {
		log.Println("Error checking if book already exists in cart:", err)
		return 2, err
	}
	if count > 0 {
		fmt.Println("Book already exists in the cart")
		return 1, nil
	}
	// Get book information from the books table
	var bookName, author string
	var cost float64
	err = db.QueryRow("SELECT book_name, author, cost FROM books WHERE bookid = $1", bookid).Scan(&bookName, &author, &cost)
	if err != nil {
		log.Println("Error retrieving book information:", err)
		return 2, err
	}

	// Insert the book into the cart
	_, err = db.Exec("INSERT INTO cart (uid, bookid, book_name, author, cost) VALUES ($1, $2, $3, $4, $5)",
		uid, bookid, bookName, author, cost)
	if err != nil {
		log.Println("Error adding book to cart:", err)
		return 2, err
	}
	fmt.Println("Book added to cart successfully")
	return 0, nil
}

// 0 bought
// 1 already owned books not bought
func buyBooks(db *sql.DB, uid int) (int, error, []Book) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error during transaction: &s", err)
		return 1, err, nil
	}
	defer tx.Rollback()

	// Get all books from the cart for the user
	rows, err := tx.Query("SELECT c.bookid, b.book_name, b.genre, b.download_url,b.author FROM cart c JOIN books b ON c.bookid = b.bookid WHERE c.uid = $1", uid)
	if err != nil {
		log.Println("Error during transaction: &s", err)
		return 1, err, nil
	}
	defer rows.Close()
	var recc []Book
	for rows.Next() {
		var bookID int
		var bookName, genre, download_url, author string
		if err := rows.Scan(&bookID, &bookName, &genre, &download_url, &author); err != nil {
			log.Println("Error during transaction: &s", err)
			return 1, err, nil
		}
		rows.Close() //i added

		recc, err = getBooksByGenreOrAuthor(db, genre, author)
		if err != nil {
			log.Println("Error during transaction: &s", err)
			return 1, err, nil
		}

		// Check if the book has already been bought
		var isBought bool
		err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM bought_books WHERE uid = $1 AND bookid = $2)", uid, bookID).Scan(&isBought)
		if err != nil {
			log.Println("Error during transaction:", err)
			return 1, err, nil
		}
		if isBought {
			fmt.Printf("Book %d is already bought, please remove it from the cart to buy other books\n", bookID)
			log.Printf("Book %s is already bought\n", bookName)
			return 1, nil, nil
		}

		// Move the book from cart to bought_books
		_, err = tx.Exec("INSERT INTO bought_books (uid, bookid, book_name, genre, download_url) VALUES ($1, $2, $3, $4, $5)",
			uid, bookID, bookName, genre, download_url)
		if err != nil {
			log.Println("Error during transaction: &s", err)
			return 1, err, nil
		}

		// Remove the book from the cart
		_, err = tx.Exec("DELETE FROM cart WHERE uid = $1 AND bookid = $2", uid, bookID)
		if err != nil {
			log.Println("Error during transaction: &s", err)
			return 1, err, nil
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error during transaction: &s", err)
		return 1, err, nil
	}

	fmt.Println("Books bought successfully, if present in cart")
	return 0, nil, recc
}

func deleteFromCart(db *sql.DB, uid, bookid int) {
	_, err := db.Exec("DELETE FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid)
	if err != nil {
		log.Println("Error deleting book from cart:", err)
		panic(err)
	}
	fmt.Println("Book deleted from cart successfully, if present")
}



func viewOwnedBooks(db *sql.DB, uid int) ([]Book, error) {
	rows, err := db.Query("SELECT bookid, book_name, genre FROM bought_books WHERE uid = $1", uid)
	if err != nil {
		fmt.Println("Error retrieving owned books:", err)
		return nil, err
	}
	defer rows.Close()

	var ownedBooks []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.Genre)
		if err != nil {
			fmt.Println("Error retrieving owned books from rows:", err)
			return nil, err
		}
		ownedBooks = append(ownedBooks, book)
	}
	fmt.Println("Owned Books:")
	return ownedBooks, nil
}

// 0 success
// 1 review exists/failed
func giveReview(db *sql.DB, uid, bookID int, review string) (int, error) {
	// Check if a review by the same user for the same book already exists
	var existingReview int
	err := db.QueryRow("SELECT reviewid FROM reviews WHERE uid = $1 AND bookid = $2", uid, bookID).Scan(&existingReview)
	if err == nil {
		fmt.Println("A review by the same user for the same book already exists")
		return 1, nil
	} else if err != sql.ErrNoRows {
		log.Println("Error checking for existing review:", err)
		return 1, err
	}

	// Insert the new review
	_, err = db.Exec("INSERT INTO reviews (uid, bookid, review) VALUES ($1, $2, $3)", uid, bookID, review)
	if err != nil {
		log.Println("Error giving review:", err)
		return 1, err
	}
	fmt.Println("Review added successfully")
	return 0, nil
}

func displayBookReviews(db *sql.DB, bookID int) error {
	rows, err := db.Query("SELECT review FROM reviews WHERE bookid = $1", bookID)
	if err != nil {
		log.Println("Error retrieving reviews:", err)
		return err
	}
	defer rows.Close()

	fmt.Println("Reviews for Book ID:", bookID)
	fmt.Println("=============================")

	for rows.Next() {
		var review string
		if err := rows.Scan(&review); err != nil {
			log.Println("Error retrieving reviews:", err)
			return err
		}
		fmt.Println(review)
	}
	return nil
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

func viewCart(db *sql.DB, uid int) ([]CartItem, error) {
	rows, err := db.Query("SELECT c.bookid, b.book_name, b.author, b.genre, b.cost FROM cart c JOIN books b ON c.bookid = b.bookid WHERE c.uid = $1", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []CartItem
	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.BookID, &item.BookName, &item.Author, &item.Genre, &item.Cost); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}

	return cartItems, nil
}
func getBooksByGenreOrAuthor(db *sql.DB, genre, author string) ([]Book, error) {
	// Query the database to select books that match the genre or author
	query := "SELECT bookid, book_name, author, genre FROM books WHERE genre = $1 OR author = $2"
	rows, err := db.Query(query, genre, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Genre); err != nil {
			fmt.Print("books", books)
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
