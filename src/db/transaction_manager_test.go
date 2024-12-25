package db

import (
	"github.com/anmollp/generic-db-go/src/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

const NumWorkers = 4

func TestTransactionManagerRegistryMultipleThreads(t *testing.T) {
	var wg sync.WaitGroup
	var tmr = TxManagerPool

	verifyTx := func() {
		defer wg.Done()

		// Register a transaction manager
		tmr.Register()
		txManager := tmr.GetTransactionManager()
		time.Sleep(1 * time.Second)

		// Check the number of transaction managers and trackers
		numTransactionManagers := tmr.getNumTransactionManagers()
		if numTransactionManagers != NumWorkers {
			t.Errorf("Expected %d tx managers, got %d", NumWorkers, numTransactionManagers)
		}

		numTransactionManagersTrackers := tmr.getNumTransactionManagerTrackers()
		if numTransactionManagersTrackers != NumWorkers {
			t.Errorf("Expected %d tx managers, got %d", NumWorkers, numTransactionManagersTrackers)
		}

		// Validate the transaction manager for this goroutine
		if txManager != tmr.GetTransactionManager() {
			t.Errorf("Mismatch in tx manager for go routine %d", utils.GetGoroutineID())
		}
		txManager2 := tmr.GetTransactionManager()
		if txManager != txManager2 {
			t.Errorf("Expected the same tx manager, got different instances")
		}
		time.Sleep(1 * time.Second)

		// Release the transaction manager
		tmr.Release(true)
		time.Sleep(1 * time.Second)

		// Validate that transaction managers and trackers are cleaned up
		numTransactionManagers = tmr.getNumTransactionManagers()
		if numTransactionManagers != 0 {
			t.Errorf("Expected 0 tx managers after release, got %d", numTransactionManagers)
		}
		numTransactionManagersTrackers = tmr.getNumTransactionManagerTrackers()
		if numTransactionManagersTrackers != 0 {
			t.Errorf("Expected 0 tx trackers after release, got %d", numTransactionManagersTrackers)
		}
	}

	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		go verifyTx()
	}

	wg.Wait()
}

func TestTransactionManagerRegistryNoRegisteringException(t *testing.T) {
	// Trying to get tx manager without first registering should panic
	assert.PanicsWithValue(t,
		"No TransactionManager registered for this goroutine. Please call Register() first.",
		func() {
			TxManagerPool.GetTransactionManager()
		})
}

func TestTransactionManagerCommitLifecycle(t *testing.T) {
	var tmr = TxManagerPool
	tmr.Register()
	txManager := tmr.GetTransactionManager()

	tmr.Register()
	txManager2 := tmr.GetTransactionManager()

	if txManager != txManager2 {
		t.Errorf("Expected same transaction manager, but got different ones")
	}

	// Since tx manager was registered twice, it will not commit when release() is called only once
	tmr.Release(true)
	assert.Equal(t, tmr.getNumTransactionManagers(), 1)
	assert.Equal(t, tmr.getNumTransactionManagerTrackers(), 1)

	// It will only commit when all components using the tx manager releases it
	tmr.Release(true)
	assert.Equal(t, tmr.getNumTransactionManagers(), 0)
	assert.Equal(t, tmr.getNumTransactionManagerTrackers(), 0)
}

func TestTransactionManagerReleaseLifeCycle(t *testing.T) {
	var tmr = TxManagerPool
	tmr.Register()
	tmr.Register()

	// Rollbacks will always be done immediately even if tx manager is registered multiple times
	tmr.Release(false)
	assert.Equal(t, tmr.getNumTransactionManagers(), 0)
	assert.Equal(t, tmr.getNumTransactionManagerTrackers(), 0)
}
