package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/pkg/database"
)

func TestConnectSQLite_InMemory(t *testing.T) {
	db, err := database.ConnectSQLite(":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	var result int
	err = db.QueryRowContext(context.Background(), "SELECT 1").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestDBManager_RegisterAndRetrieve(t *testing.T) {
	mgr := database.NewManager("sqlite1")

	db1, err := database.ConnectSQLite(":memory:")
	require.NoError(t, err)
	defer db1.Close()

	db2, err := database.ConnectSQLite(":memory:")
	require.NoError(t, err)
	defer db2.Close()

	mgr.RegisterSQL("sqlite1", db1)
	mgr.RegisterSQL("sqlite2", db2)

	retrieved1, err := mgr.SQL("sqlite1")
	require.NoError(t, err)
	assert.Equal(t, db1, retrieved1)

	retrievedDefault := mgr.DefaultSQL()
	assert.Equal(t, db1, retrievedDefault)

	retrieved2, err := mgr.SQL("sqlite2")
	require.NoError(t, err)
	assert.Equal(t, db2, retrieved2)

	_, err = mgr.SQL("unknown")
	assert.Error(t, err)

	mgr.CloseAll(context.Background())
}
