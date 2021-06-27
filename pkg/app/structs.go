package app

type Image struct {
	Title       string
	Description string
	Link        string
	Video       bool
}

type Album struct {
	Title       string
	Description string
	Images      []Image
}
