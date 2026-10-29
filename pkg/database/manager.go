package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/uptrace/bun"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DBManager struct {
	mu           sync.RWMutex
	defaultName  string
	sqlConns     map[string]*bun.DB
	mongoConns   map[string]*mongo.Database
	mongoClients map[string]*mongo.Client
}

func NewManager(defaultName string) *DBManager {
	return &DBManager{
		defaultName:  defaultName,
		sqlConns:     make(map[string]*bun.DB),
		mongoConns:   make(map[string]*mongo.Database),
		mongoClients: make(map[string]*mongo.Client),
	}
}

func (m *DBManager) RegisterSQL(name string, db *bun.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sqlConns[name] = db
}

func (m *DBManager) RegisterMongo(name string, client *mongo.Client, db *mongo.Database) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mongoClients[name] = client
	m.mongoConns[name] = db
}

func (m *DBManager) SQL(name string) (*bun.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.sqlConns[name]
	if !ok {
		return nil, fmt.Errorf("sql connection '%s' not registered", name)
	}
	return db, nil
}

func (m *DBManager) DefaultSQL() *bun.DB {
	db, _ := m.SQL(m.defaultName)
	return db
}

func (m *DBManager) Mongo(name string) (*mongo.Database, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.mongoConns[name]
	if !ok {
		return nil, fmt.Errorf("mongo connection '%s' not registered", name)
	}
	return db, nil
}

// CloseAll menutup seluruh pool koneksi SQL dan mendiskoneksikan client MongoDB secara aman
func (m *DBManager) CloseAll(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, db := range m.sqlConns {
		_ = db.Close()
	}

	for _, client := range m.mongoClients {
		_ = client.Disconnect(ctx)
	}
}
