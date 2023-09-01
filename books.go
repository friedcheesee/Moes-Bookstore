package main

import (
	"database/sql"
	"fmt"
	"log"
)

// codes returned by this function, to debug/show status of cart addition
// 0 added
// 1 already exists
// 2 internal error
func addToCart(db *sql.DB, uid, bookid int) (int, error) {
	// Check if the book already exists in the cart
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid).Scan(&count)
	//check the count of statements where the uid owns bookid. if already owned,cannot buy again.
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
	logEvent("Book added to cart successfully")
	fmt.Println("Book added to cart successfully")
	return 0, nil
}

// codes returned by this function, to debug/show status if book was successfully bought
// 0 bought
// 1 already owned books not bought
func buyBooks(db *sql.DB, uid int) (int, error, []Book) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error during transaction:", err)
		return 1, err, nil
	}
	defer tx.Rollback()

	// Get all books from the cart for the user
	rows, err := tx.Query("SELECT c.bookid, b.book_name, b.genre, b.download_url,b.author FROM cart c JOIN books b ON c.bookid = b.bookid WHERE c.uid = $1", uid)
	if err != nil {
		log.Println("Error during transaction:", err)
		return 1, err, nil
	}
	defer rows.Close()

	//store in the empty book structure, which will be later returned to the user by the handler
	var recc []Book
	for rows.Next() {
		var bookID int
		var bookName, genre, download_url, author string
		if err := rows.Scan(&bookID, &bookName, &genre, &download_url, &author); err != nil {
			log.Println("Error during transaction: ", err)
			return 1, err, nil
		}
		rows.Close() 

		//to get recomended books based on the most recent purchase
		recc, err = getBooksByGenreOrAuthor(db, genre, author)
		if err != nil {
			log.Println("Error during transaction: ", err)
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
			logEvent("Book already bought "+bookName)
			return 1, nil, nil
		}

		// Move the book from cart to bought_books
		_, err = tx.Exec("INSERT INTO bought_books (uid, bookid, book_name, genre, download_url) VALUES ($1, $2, $3, $4, $5)",
			uid, bookID, bookName, genre, download_url)
		if err != nil {
			log.Println("Error during transaction: ", err)
			return 1, err, nil
		}

		// Remove the book from the cart
		_, err = tx.Exec("DELETE FROM cart WHERE uid = $1 AND bookid = $2", uid, bookID)
		if err != nil {
			log.Println("Error during transaction: ", err)
			return 1, err, nil
		}
	}

	//commit the transaction, if any errors encountered, transaction rollbacks
	if err := tx.Commit(); err != nil {
		log.Println("Error during transaction: ", err)
		return 1, err, nil
	}
	logEvent("Books bought successfully, if present in cart")
	fmt.Println("Books bought successfully, if present in cart")
	return 0, nil, recc
}

// function to delete items from user cart
func deleteFromCart(db *sql.DB, uid, bookid int) {
	_, err := db.Exec("DELETE FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid)
	if err != nil {
		log.Println("Error deleting book from cart:", err)
		panic(err)
	}
	logEvent("Book deleted from cart successfully, if present")
	fmt.Println("Book deleted from cart successfully, if present")
}

//function to see a users inventory
func viewOwnedBooks(db *sql.DB, uid int) ([]Book, error) {
	rows, err := db.Query("SELECT bookid, book_name, genre, download_url FROM bought_books WHERE uid = $1", uid)
	if err != nil {
		fmt.Println("Error retrieving owned books:", err)
		return nil, err
	}
	defer rows.Close()

	//using array of book structs to store the books
	var ownedBooks []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.Genre, &book.DownloadURL)
		if err != nil {
			fmt.Println("Error retrieving owned books from rows:", err)
			return nil, err
		}
		ownedBooks = append(ownedBooks, book)
	}
	logEvent("Owned books retrieved successfully")
	return ownedBooks, nil
}

// codes returned by this function, to debug/show status of review addition
// 0 success
// 1 review exists/failed
func giveReview(db *sql.DB, uid, bookID int, review string) (int, error) {
	
	// Check if a review by the same user for the same book already exists
	var existingReview int
	err := db.QueryRow("SELECT reviewid FROM reviews WHERE uid = $1 AND bookid = $2", uid, bookID).Scan(&existingReview)
	if err == nil {
		fmt.Println("A review by the same user for the same book already exists")
		logEvent("A review by the same user for the same book already exists")
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
	logEvent("Review added successfully")
	fmt.Println("Review added successfully")
	return 0, nil
}

// to display reviewd for a book - not used yet.
// func displayBookReviews(db *sql.DB, bookID int) error {
// 	rows, err := db.Query("SELECT review FROM reviews WHERE bookid = $1", bookID)
// 	if err != nil {
// 		log.Println("Error retrieving reviews:", err)
// 		return err
// 	}
// 	defer rows.Close()


// 	for rows.Next() {
// 		var review string
// 		if err := rows.Scan(&review); err != nil {
// 			log.Println("Error retrieving reviews:", err)
// 			return err
// 		}
// 		fmt.Println(review)
// 	}
// 	return nil
// }

// to search a book based on query(bookname), genre and author
func searchBooks(db *sql.DB, query, genre, author string) ([]Book, error) {
	//query to search books
	rows, err := db.Query("SELECT bookid, book_name, author, genre, cost FROM books WHERE book_name ILIKE $1 AND genre ILIKE $2 AND author ILIKE $3",
		"%"+query+"%", "%"+genre+"%", "%"+author+"%")
	if err != nil {
		CheckError(err)
		return nil, err
	}
	defer rows.Close()

	//storing results in an array of book structs
	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Genre, &book.Cost)
		if err != nil {
			return nil, err
		}
		book.DownloadURL="buy the book to access URL"
		books = append(books, book)
	}
	logEvent("Books retrieved successfully")
	return books, nil
}

//function to see cart items of a user
func viewCart(db *sql.DB, uid int) ([]CartItem, error) {
	// Query the database to select books that match the genre or author
	rows, err := db.Query("SELECT c.bookid, b.book_name, b.author, b.genre, b.cost FROM cart c JOIN books b ON c.bookid = b.bookid WHERE c.uid = $1", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//store in the empty cartitem structure, which will be later returned to the user by the handler
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
//reccomendation system, selects relevant book based on recent purchase
func getBooksByGenreOrAuthor(db *sql.DB, genre, author string) ([]Book, error) {
	// Query the database to select books that match the genre or author
	query := "SELECT bookid, book_name, author, genre FROM books WHERE genre = $1 OR author = $2"
	rows, err := db.Query(query, genre, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//store in the empty book structure, which will be later returned to the user by the handler
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
