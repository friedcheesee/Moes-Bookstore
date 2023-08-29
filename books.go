package main
import (
	"fmt"
"log"
	"database/sql"
)
func addToCart(db *sql.DB, uid, bookid int){
	// Check if the book already exists in the cart
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cart WHERE uid = $1 AND bookid = $2", uid, bookid).Scan(&count)
	//check the count of statements where the uid owns bookid. if owns 1 time, cant own again.
	CheckError(err)
	if count > 0 {
		fmt.Println("Book already exists in the cart")
		return
	}
	// Get book information from the books table
	var bookName, author string
	var cost float64
	err = db.QueryRow("SELECT book_name, author, cost FROM books WHERE bookid = $1", bookid).Scan(&bookName, &author, &cost)
	CheckError(err)

	// Insert the book into the cart
	_, err = db.Exec("INSERT INTO cart (uid, bookid, book_name, author, cost) VALUES ($1, $2, $3, $4, $5)",
		uid, bookid, bookName, author, cost)
	CheckError(err)
	fmt.Println("Book added to cart successfully")
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
    rows, err := db.Query("SELECT bookid, book_name, author, genre, cost FROM books")
    if err != nil {
        log.Println("Error retrieving available books:", err)
        return err
    }
    defer rows.Close()

    fmt.Println("Available Books:")
    fmt.Println("================")
    for rows.Next() {
        var bookID int
        var bookName, author, genre string
        var cost float64
        if err := rows.Scan(&bookID, &bookName, &author, &genre, &cost); err != nil {
            log.Println("Error retrieving available books:", err)
            return err
        }
        fmt.Printf("Book ID: %d\nTitle: %s\nAuthor: %s\nGenre: %s\nCost: $%.2f\n\n",
            bookID, bookName, author, genre, cost)
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