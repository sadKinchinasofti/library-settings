package handlers

type Library struct {
	BookID     int //`json:"id"`
	BookName   string
	BookAuthor string
	//BookQuantity int
}

type GetLibrary struct {
	BookID     string
	BookName   string
	BookAuthor string
}
