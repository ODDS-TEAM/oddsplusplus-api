package model

// Book ...
type Book struct {
	BookName string `json:"book_name"`
	Author   string `json:"author"`
	Imgage   string `json:"img"`
	Price    string `json:"price"`
	Format    string `json:"format" bson:"format"`
}

// BookUrl ...
type BookUrl struct {
	BookUrl string `json:"book_url"`
}
type FindError struct {
	Message string `json:"message"`
}
