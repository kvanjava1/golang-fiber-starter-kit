package mysql1

import (
	"context"

	"github.com/uptrace/bun"

	"golang-starter-kit/database/seeders"
)

func init() {
	seeders.Register("mysql1", NewExampleMySQLSeeder())
}

type ExampleMySQLSeeder struct{}

func NewExampleMySQLSeeder() *ExampleMySQLSeeder {
	return &ExampleMySQLSeeder{}
}

func (s *ExampleMySQLSeeder) Name() string {
	return "mysql1:ExampleSeeder"
}

func (s *ExampleMySQLSeeder) Run(ctx context.Context, db bun.IDB) error {
	// Masukkan kode seeder MySQL Anda di sini.
	// Query ini HANYA akan dieksekusi saat 'go run cmd/make/main.go db:seed --db=mysql1' dijalankan.
	return nil
}
