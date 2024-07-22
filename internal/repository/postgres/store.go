package postgres

import (
    "context"
    "database/sql"
)

// Store provides methods to interact with the database, supporting transactions.
type Store struct {
    *Queries
    db *sql.DB
}

// NewStore creates a new Store instance.
func NewStore(db *sql.DB) *Store {
    return &Store{
        Queries: New(db),
        db:      db,
    }
}

// BeginTx starts a new transaction.
func (s *Store) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
    tx, err := s.db.BeginTx(ctx, opts)
    if err != nil {
        return nil, err
    }
    return &Tx{
        Tx:      tx,
        Queries: s.Queries,
    }, nil
}

// Tx wraps an *sql.Tx and provides methods for interacting with the database within a transaction.
type Tx struct {
    *sql.Tx
    *Queries
}

// CreateOrder creates a new order within the transaction.
func (tx *Tx) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
    return tx.Queries.CreateOrder(ctx, arg)
}

// DeleteOrder deletes an order within the transaction.
func (tx *Tx) DeleteOrder(ctx context.Context, id int64) error {
    return tx.Queries.DeleteOrder(ctx, id)
}

// GetOrder retrieves an order by ID within the transaction.
func (tx *Tx) GetOrder(ctx context.Context, id int64) (Order, error) {
    return tx.Queries.GetOrder(ctx, id)
}

// ListOrders retrieves all orders within the transaction.
func (tx *Tx) ListOrders(ctx context.Context) ([]Order, error) {
    return tx.Queries.ListOrders(ctx)
}

// SearchOrdersByStatus searches orders by status within the transaction.
func (tx *Tx) SearchOrdersByStatus(ctx context.Context, status OrderStatus) ([]Order, error) {
    return tx.Queries.SearchOrdersByStatus(ctx, status)
}

// SearchOrdersByUser searches orders by user ID within the transaction.
func (tx *Tx) SearchOrdersByUser(ctx context.Context, userID int64) ([]Order, error) {
    return tx.Queries.SearchOrdersByUser(ctx, userID)
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
    return tx.Tx.Commit()
}

// Rollback rolls back the transaction.
func (tx *Tx) Rollback() error {
    return tx.Tx.Rollback()
}
