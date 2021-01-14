package models

type Product struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// gorm.Model
// ID    uint   `json:"id"`
// Name string `json:"name"`
// Desc  string `json:"desc"`
// Price string `json:"price"`
// Image string `json:"image"`
