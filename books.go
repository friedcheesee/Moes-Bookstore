package main
import (
	"fmt"
	"log"
	"database/sql"
)

// 0 added
// 1 already exists
// 2 internal error
func addToCart(db *sql.DB, uid, bookid int) (int ,error) {
	// Check if the book already exists in the cart
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid).Scan(&count)
	//check the count of statements where the uid owns bookid. if owns 1 time, cant own again.
	if err != nil {
		log.Println("Error checking if book already exists in cart:", err)
		return 2,err
	}
	if count > 0 {
		fmt.Println("Book already exists in the cart")
		return 1,nil
	}
	// Get book information from the books table
	var bookName, author string
	var cost float64
	err = db.QueryRow("SELECT book_name, author, cost FROM books WHERE bookid = $1", bookid).Scan(&bookName, &author, &cost)
	if err != nil {
		log.Println("Error retrieving book information:", err)
		return 2,err
	}

	// Insert the book into the cart
	_, err = db.Exec("INSERT INTO cart (uid, bookid, book_name, author, cost) VALUES ($1, $2, $3, $4, $5)",
		uid, bookid, bookName, author, cost)
	if err != nil {
		log.Println("Error adding book to cart:", err)
		return 2,err
	} 
	fmt.Println("Book added to cart successfully")
	return 0,nil
}

func buyBooks(db *sql.DB, uid int) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error during transaction: &s", err)
		return err
	}
	defer tx.Rollback()

	// Get all books from the cart for the user
	rows, err := tx.Query("SELECT c.bookid, b.book_name, b.genre, b.download_url FROM cart c JOIN books b ON c.bookid = b.bookid WHERE c.uid = $1", uid)
	if err != nil {
		log.Println("Error during transaction: &s", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var bookID int
		var bookName, genre, download_url string
		if err := rows.Scan(&bookID, &bookName, &genre, &download_url); err != nil {
			log.Println("Error during transaction: &s", err)
			return err
		}
		rows.Close() //i added

		// Check if the book has already been bought
		var isBought bool
		err := tx.QueryRow("SELECT EXISTS (SELECT 1 FROM bought_books WHERE uid = $1 AND bookid = $2)", uid, bookID).Scan(&isBought)
		if err != nil {
			log.Println("Error during transaction:", err)
			return err
		}
		if isBought {
			fmt.Printf("Book %s is already bought, please remove it from the cart to buy other books\n", bookName)
			log.Printf("Book %s is already bought\n", bookName)
			return nil
		}

		// Move the book from cart to bought_books
		_, err = tx.Exec("INSERT INTO bought_books (uid, bookid, book_name, genre, download_url) VALUES ($1, $2, $3, $4, $5)",
		uid, bookID, bookName, genre, download_url)
		if err != nil {
			log.Println("Error during transaction: &s", err)
			return err
		}

		// Remove the book from the cart
		_, err = tx.Exec("DELETE FROM cart WHERE uid = $1 AND bookid = $2", uid, bookID)
		if err != nil {
			log.Println("Error during transaction: &s", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error during transaction: &s", err)
		return err
	}

	fmt.Println("Books bought successfully")
	return nil
}

func deleteFromCart(db *sql.DB, uid, bookid int) {
    _, err := db.Exec("DELETE FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid)
    if err != nil {
        log.Println("Error deleting book from cart:", err)
        panic(err)
    }
	fmt.Println("Book deleted from cart successfully")
}

func displayAvailableBooks(db *sql.DB) error {
    rows, err := db.Query("SELECT b.bookid, b.book_name, b.author, b.genre, b.cost, r.review "+
        "FROM books b LEFT JOIN reviews r ON b.bookid = r.bookid")
    if err != nil {
        log.Println("Error retrieving available books:", err)
        return err
    }
    defer rows.Close()

    fmt.Println("Available Books:")
    fmt.Println("================")
    
    var currentBookID int
    var currentBookName, currentAuthor, currentGenre string
    var currentCost float64
    
    for rows.Next() {
        var bookID int
        var bookName, author, genre string
        var cost float64
        var review sql.NullString
        if err := rows.Scan(&bookID, &bookName, &author, &genre, &cost, &review); err != nil {
            log.Println("Error retrieving available books:", err)
            return err
        }
        
        if currentBookID != bookID {
            if currentBookID != 0 {
                fmt.Println() // Print a new line before the next book
            }
            currentBookID = bookID
            currentBookName = bookName
            currentAuthor = author
            currentGenre = genre
            currentCost = cost

            fmt.Printf("Book ID: %d\nTitle: %s\nAuthor: %s\nGenre: %s\nCost: $%.2f\n",
                currentBookID, currentBookName, currentAuthor, currentGenre, currentCost)
        }

        if review.Valid {
            fmt.Printf("Review: %s\n", review.String)
        } else {
            fmt.Println("No reviews available")
        }
    }
    return nil
}


func viewOwnedBooks(db *sql.DB, uid int) error {
	rows, err := db.Query("SELECT bookid, book_name, genre FROM bought_books WHERE uid = $1", uid)
	if err != nil {
		log.Println("Error retrieving owned books:", err)
		return err
	}
	defer rows.Close()

	fmt.Println("Owned Books:")
	fmt.Println("=============")
	for rows.Next() {
		var bookID int
		var bookName, genre string
		if err := rows.Scan(&bookID, &bookName, &genre); err != nil {
			log.Println("Error retrieving owned books:", err)
			return err
		}
		fmt.Printf("Book ID: %d\nTitle: %s\nGenre: %s\n\n", bookID, bookName, genre)
	}
	return nil
}

func giveReview(db *sql.DB, uid, bookID int, review string) error {
	// Check if a review by the same user for the same book already exists
	var existingReview int
	err := db.QueryRow("SELECT reviewid FROM reviews WHERE uid = $1 AND bookid = $2", uid, bookID).Scan(&existingReview)
	if err == nil {
		fmt.Println("A review by the same user for the same book already exists")
		return nil
	} else if err != sql.ErrNoRows {
		log.Println("Error checking for existing review:", err)
		return err
	}

	// Insert the new review
	_, err = db.Exec("INSERT INTO reviews (uid, bookid, review) VALUES ($1, $2, $3)", uid, bookID, review)
	if err != nil {
		log.Println("Error giving review:", err)
		return err
	}
	fmt.Println("Review added successfully")
	return nil
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
