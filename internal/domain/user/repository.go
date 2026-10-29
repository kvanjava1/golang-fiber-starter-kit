package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"golang-starter-kit/pkg/response"
)

type Repository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, item *User) error
	Update(ctx context.Context, item *User) error
	Delete(ctx context.Context, id string) error
}

type bunRepository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &bunRepository{db: db}
}

func (r *bunRepository) FindAll(ctx context.Context) ([]User, error) {
	var items []User
	err := r.db.NewSelect().
		Model(&items).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return items, nil
}

func (r *bunRepository) FindByID(ctx context.Context, id string) (*User, error) {
	item := new(User)
	err := r.db.NewSelect().
		Model(item).
		Where("u.id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("failed to query user by id %s: %w", id, err)
	}
	return item, nil
}

func (r *bunRepository) Create(ctx context.Context, item *User) error {
	_, err := r.db.NewInsert().
		Model(item).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *bunRepository) Update(ctx context.Context, item *User) error {
	res, err := r.db.NewUpdate().
		Model(item).
		WherePK().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *bunRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.NewDelete().
		Model((*User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return response.ErrNotFound
	}
	return nil
}
