package database

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ConnectMongo membuat koneksi client dan database MongoDB
func ConnectMongo(ctx context.Context, uri string) (*mongo.Client, *mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	// Ambil database name dari URI (misal mongodb://localhost:27017/my_db)
	dbName := "app"
	parts := strings.Split(uri, "/")
	if len(parts) > 3 {
		lastPart := parts[len(parts)-1]
		cleanName := strings.Split(lastPart, "?")[0]
		if cleanName != "" {
			dbName = cleanName
		}
	}

	return client, client.Database(dbName), nil
}
