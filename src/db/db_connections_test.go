package db_test

import (
	"github.com/anmollp/generic-db-go/src/db"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"testing"
)

func TestMultiThreadedConnection(t *testing.T) {

	const numWorkers = 2
	var results []int
	var wg sync.WaitGroup

	// Function to run in each goroutine
	queryFunc := func(x int, wg *sync.WaitGroup, results *[]int) {
		defer wg.Done()
		result, _ := db.RdsConnPool.ExecuteQuery("select ? AS number;", []interface{}{x}, true)
		*results = append(*results, int(result.(map[string]interface{})["number"].(int64)))
	}

	// Using WaitGroup to synchronize goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go queryFunc(i, &wg, &results)
	}

	wg.Wait() // Wait for all goroutines to finish

	if len(results) != numWorkers {
		t.Errorf("Expected %d, got %d", numWorkers, len(results))
	}
}
