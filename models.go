package main

type Book struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Author      string  `json:"author"`
	Genre       string  `json:"genre"`
	Cost        float64 `json:"cost"`
	DownloadURL string  `json:"download_url"`
}
type CartItem struct {
	BookID   int     `json:"bookid"`
	BookName string  `json:"book_name"`
	Author   string  `json:"author"`
	Genre    string  `json:"genre"`
	Cost     float64 `json:"cost"`
}

type BuyBooksResponse struct {
	Code           int
	Recommendation string
}

type User struct {
	ID       int
	Username string
	Email    string
}
