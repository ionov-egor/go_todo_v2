package db

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
