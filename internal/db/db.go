package db

import (
	"context"
	"fmt"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	Task   string
	Status string
}

type DB struct {
	pool *pgxpool.Pool
}

func New(user, password, dbname, host string, port int) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the db: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	return &DB{pool: pool}, nil
}

func (db *DB) InsertItem(ctx context.Context, item Item) error {
	query := `INSERT INTO todo_items (task, status) VALUES ($1, $2)`
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status)
	return err
}

func (db *DB) GetAllItems(ctx context.Context) ([]Item, error) {
	query := `SELECT task, status FROM todo_items`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Task, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, rows.Err()
	}
	return items, nil
}

func (db *DB) DeleteItem(ctx context.Context, id int) error {
	query := `DELETE FROM todo_items WHERE id = $1`
	_, err := db.pool.Exec(ctx, query, id)
	return err
}

func (db *DB) Close() {
	db.pool.Close()
}
