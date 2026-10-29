package seeders_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"golang-starter-kit/database/seeders"
	"golang-starter-kit/pkg/database"
)

type dummySeederSuccess struct{}

func (d *dummySeederSuccess) Name() string { return "dummySuccess" }
func (d *dummySeederSuccess) Run(ctx context.Context, db bun.IDB) error {
	_, err := db.ExecContext(ctx, "INSERT INTO items (name) VALUES (?)", "Item 1")
	return err
}

type dummySeederFail struct{}

func (d *dummySeederFail) Name() string { return "dummyFail" }
func (d *dummySeederFail) Run(ctx context.Context, db bun.IDB) error {
	_, _ = db.ExecContext(ctx, "INSERT INTO items (name) VALUES (?)", "Item Fail")
	return errors.New("simulated error")
}

func TestSeederRunner_SuccessAndRollback(t *testing.T) {
	ctx := context.Background()
	db, err := database.ConnectSQLite(":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Buat tabel pengujian
	_, err = db.ExecContext(ctx, "CREATE TABLE items (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	require.NoError(t, err)

	// 1. Sukses Run
	runner := seeders.NewRunner()
	runner.Register(&dummySeederSuccess{})

	err = runner.RunAll(ctx, db)
	require.NoError(t, err)

	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM items").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// 2. Gagal Run -> Harus Rollback otomatis!
	failRunner := seeders.NewRunner()
	failRunner.Register(&dummySeederFail{})

	err = failRunner.RunAll(ctx, db)
	require.Error(t, err)

	// Verifikasi count tidak bertambah menjadi 2 karena telah di-rollback
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM items").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
