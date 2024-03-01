package db

type Product struct {
	Name       string
	Mark       float32
	Categories []Category
}

type Category struct {
	Name string
}
