package app

type Image struct {
	Title       string
	Description string
	Link        string
	Video       bool
	Width       float64
	Height      float64
}

type Album struct {
	Title       string
	Description string
	Images      []Image
}
