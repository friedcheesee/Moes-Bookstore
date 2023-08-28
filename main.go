package main

import (
	//"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "joe"
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
	username := "friedcheese"
	password := "abcd"
	reguser(db, username, password)
	//logindb(db, username, password)
	ff := getuserid(db, username)
	readID(ff)
	reactivate(db, username, password)
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