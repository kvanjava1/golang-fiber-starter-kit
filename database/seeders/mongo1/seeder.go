package mongo1

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"golang-starter-kit/database/seeders"
)

func init() {
	seeders.RegisterMongo("mongo1", NewExampleMongoSeeder())
}

type ExampleMongoSeeder struct{}

func NewExampleMongoSeeder() *ExampleMongoSeeder {
	return &ExampleMongoSeeder{}
}

func (s *ExampleMongoSeeder) Name() string {
	return "mongo1:ExampleSeeder"
}

func (s *ExampleMongoSeeder) Run(ctx context.Context, db *mongo.Database) error {
	// Masukkan kode seeder MongoDB Anda di sini menggunakan db.Collection("...").
	// Kode ini HANYA akan dieksekusi saat 'go run cmd/make/main.go db:seed --db=mongo1' dijalankan.
	return nil
}
