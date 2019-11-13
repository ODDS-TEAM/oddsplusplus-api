package model

// Book ...
type Book struct {
	BookName string `json:"book_name"`
	Author   string `json:"author"`
	Imgage   string `json:"img"`
	Price    string `json:"price"`
}

// BookUrl ...
type BookUrl struct {
	BookUrl string `json:"book_url"`
}
type FindError struct {
	Message string `json:"message"`
}
