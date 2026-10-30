package sqlite1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"golang-starter-kit/database/seeders"
)

func init() {
	seeders.Register("sqlite1", NewUserSeeder())
}

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (s *UserSeeder) Name() string {
	return "sqlite1:UserSeeder"
}

func (s *UserSeeder) Run(ctx context.Context, db bun.IDB) error {
	type userRecord struct {
		bun.BaseModel `bun:"table:users"`
		ID            string    `bun:"id,pk"`
		Name          string    `bun:"name,notnull"`
		Status        string    `bun:"status,notnull"`
		CreatedAt     time.Time `bun:"created_at,notnull"`
		UpdatedAt     time.Time `bun:"updated_at,notnull"`
	}

	users := []userRecord{
		{
			ID:        uuid.New().String(),
			Name:      "Administrator (SQLite)",
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Demo User (SQLite)",
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	_, err := db.NewInsert().Model(&users).On("CONFLICT (id) DO NOTHING").Exec(ctx)
	return err
}
