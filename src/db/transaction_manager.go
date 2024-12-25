package db

import (
	"database/sql"
	"fmt"
	"sync"
)

// TransactionManager manages database connections and transactions.
type TransactionManager struct {
	connectionPool *sql.DB
	conn           *sql.Conn
	mu             sync.Mutex
}

// NewTransactionManager creates a new TransactionManager with the given connection pool.
func NewTransactionManager(pool *sql.DB) *TransactionManager {
	return &TransactionManager{
		connectionPool: pool,
	}
}

// GetConnection returns a connection from the pool, creating one if necessary.
func (tm *TransactionManager) GetConnection() (*sql.Conn, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.conn == nil {
		conn, err := tm.connectionPool.Conn(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection: %w", err)
		}
		tm.conn = conn
	}
	return tm.conn, nil
}

// Commit commits the current transaction and closes the connection.
func (tm *TransactionManager) Commit() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.conn == nil {
		return nil
	}
	defer tm.closeConnection()

	tx, err := tm.conn.BeginTx(nil, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the current transaction and closes the connection.
func (tm *TransactionManager) Rollback() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.conn == nil {
		return nil
	}
	defer tm.closeConnection()

	tx, err := tm.conn.BeginTx(nil, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("failed to roll back transaction: %w", err)
	}
	return nil
}

// closeConnection closes the current connection.
func (tm *TransactionManager) closeConnection() {
	if tm.conn != nil {
		_ = tm.conn.Close()
		tm.conn = nil
	}
}
