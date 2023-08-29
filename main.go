package main

import (
	//"fmt"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "moe"
)

func main() {
	logFile := initiatelog()
	defer logFile.Close()

	//conecct
	db := adminconnect()
	defer db.Close()
	//rows, err := getdata(db)
	//CheckError(err)
	//printnames(rows)
	//reguser(db)
	username := "banana"
	email:= "friedd@mail.com"
	password := "abcd"
	reguser(db, email, password, username)
	logindb(db, email, password)
	ff := getuserid(db, email)
	fmt.Println(ff)
	//deactivate(db, username, password)
	//addToCart(db, ff, 1)
	viewOwnedBooks(db, ff)
	//displayAvailableBooks(db)
	//deleteFromCart(db, ff, 1)
	//err := buyBooks(db, ff)
	//CheckError(err)
	if(!isUserActive(db,ff)){
		fmt.Println("User is not active")
	}
}
func CheckError(err error) {
	if err != nil {
		log.Println("Error: &s", err)
		panic(err)
	}
}
func initiatelog() *os.File{
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	CheckError(err)
	log.SetOutput(logFile)
	return logFile
}