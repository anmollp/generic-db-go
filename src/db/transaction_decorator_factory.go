package db

import "log"

type TransactionDecoratorFactory struct {
	txManagerPool *TransactionManagerRegistry
}

func NewTransactionDecoratorFactory(txManagerPool *TransactionManagerRegistry) *TransactionDecoratorFactory {
	return &TransactionDecoratorFactory{txManagerPool: txManagerPool}
}

func (factory *TransactionDecoratorFactory) Create(f func() error) func() error {
	return func() error {
		// Register a transaction manager.
		factory.txManagerPool.Register()

		// Handle transaction lifecycle.
		defer func() {
			if r := recover(); r != nil {
				// Rollback in case of a panic.
				err := factory.txManagerPool.Release(false)
				if err != nil {
					log.Printf("Failed to rollback transaction: %v", err)
				}
				panic(r) // Re-throw the panic.
			}
		}()

		// Execute the decorated function.
		err := f()
		if err != nil {
			// Rollback on error.
			releaseErr := factory.txManagerPool.Release(false)
			if releaseErr != nil {
				log.Printf("Failed to rollback transaction: %v", releaseErr)
			}
			return err
		}

		// Commit on success.
		releaseErr := factory.txManagerPool.Release(true)
		if releaseErr != nil {
			log.Printf("Failed to commit transaction: %v", releaseErr)
			return releaseErr
		}
		return nil
	}
}
