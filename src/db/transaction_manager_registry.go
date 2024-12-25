package db

import (
	"database/sql"
	"github.com/anmollp/generic-db-go/src/utils"
	"sync"
)

// TransactionManagerRegistry manages a TransactionManager for each goroutine.
type TransactionManagerRegistry struct {
	connectionPool *sql.DB
	txManagers     sync.Map // Thread-safe map to store TransactionManagers
	txTrackers     sync.Map // Thread-safe map to store usage counters
}

// Helper methods for managing the tracker count.
func (r *TransactionManagerRegistry) incrementTracker(goroutineID int) {
	count := r.getTrackerCount(goroutineID)
	r.txTrackers.Store(goroutineID, count+1)
}

func (r *TransactionManagerRegistry) decrementTracker(goroutineID int) {
	count := r.getTrackerCount(goroutineID)
	if count > 0 {
		r.txTrackers.Store(goroutineID, count-1)
	}
}

func (r *TransactionManagerRegistry) getTrackerCount(goroutineID int) int {
	count, ok := r.txTrackers.Load(goroutineID)
	if !ok {
		return 0
	}
	return count.(int)
}

// IsRegistered checks if a TransactionManager is registered for the current goroutine.
func (r *TransactionManagerRegistry) IsRegistered() bool {
	goroutineID := utils.GetGoroutineID()
	_, exists := r.txManagers.Load(goroutineID)
	return exists
}

// NewTransactionManagerRegistry creates a new TransactionManagerRegistry.
func NewTransactionManagerRegistry(pool *sql.DB) *TransactionManagerRegistry {
	return &TransactionManagerRegistry{
		connectionPool: pool,
	}
}

// Register registers a TransactionManager for the current goroutine.
func (r *TransactionManagerRegistry) Register() {
	goroutineID := utils.GetGoroutineID()
	_, managerExists := r.txManagers.Load(goroutineID)
	if !managerExists {
		r.txManagers.Store(goroutineID, NewTransactionManager(r.connectionPool))
		r.txTrackers.Store(goroutineID, 0)
	}
	r.incrementTracker(goroutineID)
}

// Release releases the TransactionManager for the current goroutine.
// Commits the transaction if `commit` is true, otherwise rolls it back.
func (r *TransactionManagerRegistry) Release(commit bool) error {
	goroutineID := utils.GetGoroutineID()

	manager, ok := r.txManagers.Load(goroutineID)
	if !ok {
		return nil // No manager registered for this goroutine.
	}
	txManager := manager.(*TransactionManager)

	if !commit {
		r.txManagers.Delete(goroutineID)
		r.txTrackers.Delete(goroutineID)
		return txManager.Rollback()
	}

	// Decrement tracker and commit only when all components have released.
	r.decrementTracker(goroutineID)
	if r.getTrackerCount(goroutineID) == 0 {
		r.txManagers.Delete(goroutineID)
		r.txTrackers.Delete(goroutineID)
		return txManager.Commit()
	}
	return nil
}

// GetTransactionManager retrieves the TransactionManager for the current goroutine.
// Panics if no TransactionManager is registered.
func (r *TransactionManagerRegistry) GetTransactionManager() *TransactionManager {
	goroutineID := utils.GetGoroutineID()
	manager, ok := r.txManagers.Load(goroutineID)
	if !ok {
		panic("No TransactionManager registered for this goroutine. Please call Register() first.")
	}
	return manager.(*TransactionManager)
}

func (r *TransactionManagerRegistry) getNumTransactionManagers() int {
	length := 0
	r.txManagers.Range(func(key, value interface{}) bool {
		length++
		return true // continue iteration
	})
	return length
}

func (r *TransactionManagerRegistry) getNumTransactionManagerTrackers() int {
	length := 0
	r.txTrackers.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}
