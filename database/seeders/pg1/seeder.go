package pg1

import (
	"context"

	"github.com/uptrace/bun"

	"golang-starter-kit/database/seeders"
)

func init() {
	seeders.Register("pg1", NewExamplePostgresSeeder())
}

type ExamplePostgresSeeder struct{}

func NewExamplePostgresSeeder() *ExamplePostgresSeeder {
	return &ExamplePostgresSeeder{}
}

func (s *ExamplePostgresSeeder) Name() string {
	return "pg1:ExampleSeeder"
}

func (s *ExamplePostgresSeeder) Run(ctx context.Context, db bun.IDB) error {
	// Masukkan kode seeder PostgreSQL Anda di sini.
	// Query ini HANYA akan dieksekusi saat 'go run cmd/make/main.go db:seed --db=pg1' dijalankan.
	return nil
}
