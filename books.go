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
			fmt.Printf("Book %s is already bought\n", bookName)
			continue
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