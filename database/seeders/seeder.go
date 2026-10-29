package seeders

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun"
)

// Seeder mendefinisikan interface seeder modular yang menerima bun.IDB (kompatibel dengan *bun.DB maupun bun.Tx)
type Seeder interface {
	Name() string
	Run(ctx context.Context, db bun.IDB) error
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
