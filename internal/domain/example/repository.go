package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"golang-starter-kit/pkg/response"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Crud, error)
	FindByID(ctx context.Context, id string) (*Crud, error)
	Create(ctx context.Context, item *Crud) error
	Update(ctx context.Context, item *Crud) error
	Delete(ctx context.Context, id string) error
}

type bunRepository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &bunRepository{db: db}
}

func (r *bunRepository) FindAll(ctx context.Context) ([]Crud, error) {
	var items []Crud
	err := r.db.NewSelect().
		Model(&items).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query cruds: %w", err)
	}
	return items, nil
}

func (r *bunRepository) FindByID(ctx context.Context, id string) (*Crud, error) {
	item := new(Crud)
	err := r.db.NewSelect().
		Model(item).
		Where("c.id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("failed to query crud by id %s: %w", id, err)
	}
	return item, nil
}

func (r *bunRepository) Create(ctx context.Context, item *Crud) error {
	_, err := r.db.NewInsert().
		Model(item).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert crud: %w", err)
	}
	return nil
}

func (r *bunRepository) Update(ctx context.Context, item *Crud) error {
	res, err := r.db.NewUpdate().
		Model(item).
		WherePK().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update crud: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *bunRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.NewDelete().
		Model((*Crud)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete crud: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return response.ErrNotFound
	}
	return nil
}
