package db

type Project struct {
	ID          int64
	Name        string
	Description string
	Resources   []*Resource
}

type Resource struct {
	ID     int64
	Name   string
	Amount int
	Goal   int
}