package models



type Product struct{
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}