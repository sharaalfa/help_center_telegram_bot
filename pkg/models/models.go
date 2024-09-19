package models

type Ticket struct {
	ID          int64
	Department  string
	Title       string
	Description string
	ClientID    int64
}
