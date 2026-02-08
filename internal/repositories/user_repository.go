package repositories

import (
	"context"
	"go-api-template/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	query, args, err := sqlx.Named(`
		INSERT INTO users (first_name, last_name, email, password)
		VALUES (:first_name, :last_name, :email, :password)
		RETURNING *`, user)
	if err != nil {
		return model.User{}, err
	}
	query = r.db.Rebind(query)

	var created model.User
	if err := r.db.GetContext(ctx, &created, query, args...); err != nil {
		return model.User{}, err
	}
	return created, nil
}

func (r *UserRepository) CheckIfUserEmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE email = $1", email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
