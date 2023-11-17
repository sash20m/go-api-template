package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Initializes the postgres driver
	"github.com/sash20m/go-api-template/internal/model"
)

type StorageInterface interface {
	AddBook(ctx context.Context, book model.AddBookRequest) (int, error)
	GetBook(ctx context.Context, id int) (model.Book, error)
	GetBooks(ctx context.Context) ([]model.Book, error)
	UpdateBook(ctx context.Context, book model.UpdateBookRequest) (int, error)
	DeleteBook(ctx context.Context, id int) error
	VerifyBookExists(ctx context.Context, id int) (bool, error)
}

// Storage contains an SQL db. Storage implements the StorageInterface.
type Storage struct {
	db *sqlx.DB
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetDB() *sqlx.DB {
	return s.db
}
