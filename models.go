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
}
