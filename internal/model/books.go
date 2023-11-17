package model

import "database/sql"

// Book is the db schema for the books table
type Book struct {
	ID       int            `json:"id" db:"id"`
	Title    string         `json:"title" db:"title"`
	Author   string         `json:"author" db:"author"`
	CoverURL string         `json:"coverUrl" db:"cover_url"`
	PostURL  sql.NullString `json:"postUrl" db:"post_url"`
	Base
}

// Below are the structures of the request/response structs in the books handler.

type AddBookRequest struct {
	Title    string `json:"title" validate:"required"`
	Author   string `json:"author" validate:"required"`
	CoverURL string `json:"coverUrl" validate:"required"`
	PostURL  string `json:"postUrl" validate:"required"`
}

type GetBookResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"coverUrl"`
	PostURL  string `json:"postUrl"`
}

type UpdateBookRequest struct {
	ID       int    `json:"id" validate:"required"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"coverUrl"`
	PostURL  string `json:"postUrl"`
}

type IDResponse struct {
	ID int `json:"id"`
}
