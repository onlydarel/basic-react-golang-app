package models

type Todo struct {
	ID     int    `json:"id"`
	Status bool   `json:"status"`
	Body   string `json:"body"`
}
