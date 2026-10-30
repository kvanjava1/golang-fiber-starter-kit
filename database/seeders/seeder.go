package seeders

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/uptrace/bun"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Seeder mendefinisikan interface seeder SQL modular yang menerima bun.IDB (kompatibel dengan *bun.DB maupun bun.Tx)
type Seeder interface {
	Name() string
	Run(ctx context.Context, db bun.IDB) error
}

// MongoSeeder mendefinisikan interface seeder NoSQL untuk MongoDB
type MongoSeeder interface {
	Name() string
	Run(ctx context.Context, db *mongo.Database) error
}

type Runner struct {
	seeders []Seeder
}

func NewRunner() *Runner {
	return &Runner{
		seeders: make([]Seeder, 0),
	}
}

func (r *Runner) Register(s Seeder) {
	r.seeders = append(r.seeders, s)
}

// RunAll mengeksekusi seluruh seeder yang terdaftar di dalam transaksi atomik
func (r *Runner) RunAll(ctx context.Context, db *bun.DB) error {
	for _, s := range r.seeders {
		slog.InfoContext(ctx, "Running seeder", "name", s.Name())
		err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			return s.Run(ctx, tx)
		})
		if err != nil {
			return fmt.Errorf("seeder '%s' failed: %w", s.Name(), err)
		}
		slog.InfoContext(ctx, "Seeder completed successfully", "name", s.Name())
	}
	return nil
}

// Global registry terisolasi per target database (sqlite1, mysql1, pg1, mongo1, dll.)
var (
	registryMu    sync.RWMutex
	sqlRegistry   = make(map[string][]Seeder)
	mongoRegistry = make(map[string][]MongoSeeder)
)

// Register mendaftarkan seeder SQL ke database target spesifik
func Register(dbTarget string, s Seeder) {
	registryMu.Lock()
	defer registryMu.Unlock()
	sqlRegistry[dbTarget] = append(sqlRegistry[dbTarget], s)
}

// RegisterMongo mendaftarkan seeder NoSQL ke database MongoDB target spesifik
func RegisterMongo(dbTarget string, s MongoSeeder) {
	registryMu.Lock()
	defer registryMu.Unlock()
	mongoRegistry[dbTarget] = append(mongoRegistry[dbTarget], s)
}

// GetSQLSeeders mengambil daftar seeder SQL HANYA untuk target database yang diminta
func GetSQLSeeders(dbTarget string) []Seeder {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return sqlRegistry[dbTarget]
}

// GetMongoSeeders mengambil daftar seeder MongoDB HANYA untuk target database yang diminta
func GetMongoSeeders(dbTarget string) []MongoSeeder {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return mongoRegistry[dbTarget]
}
