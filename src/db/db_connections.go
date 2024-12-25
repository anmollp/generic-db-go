package db

import (
	"context"
	"database/sql"
	"log"
	"sync"
)

type SQLUpdate struct {
	SQL    string
	Values [][]interface{}
}

type RDSPooledConnection struct {
	cnxPool       *sql.DB
	txManagerPool *TransactionManagerRegistry
	mu            sync.Mutex
}

func NewRDSPooledConnection(cnxPool *sql.DB, txManagerPool *TransactionManagerRegistry) *RDSPooledConnection {
	return &RDSPooledConnection{
		cnxPool:       cnxPool,
		txManagerPool: txManagerPool,
	}
}

func (r *RDSPooledConnection) ExecuteQuery(sqlQuery string, params []interface{}, fetchOne bool) (interface{}, error) {
	var cnx *sql.Conn
	var rows *sql.Rows
	var err error

	cnx, err = r.cnxPool.Conn(context.Background())
	if err != nil {
		log.Printf("Error getting connection: %v", err)
		return nil, err
	}
	defer cnx.Close()

	stmt, err := cnx.PrepareContext(context.Background(), sqlQuery)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.QueryContext(context.Background(), params...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rowData := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range rowData {
			rowPointers[i] = &rowData[i]
		}

		if err := rows.Scan(rowPointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = rowData[i]
		}
		results = append(results, rowMap)
	}

	if fetchOne && len(results) > 0 {
		return results[0], nil
	}

	return results, nil
}

func (r *RDSPooledConnection) ExecuteUpdates(updates []SQLUpdate) ([]int64, []int64, error) {
	if len(updates) == 0 {
		return []int64{}, []int64{}, nil
	}

	var cnx *sql.Conn
	var err error
	var rowCounts []int64
	var newRowIDs []int64

	txManager := r.txManagerPool.GetTransactionManager()
	if txManager == nil {
		cnx, err = r.cnxPool.Conn(context.Background())
		if err != nil {
			log.Printf("Error getting connection: %v", err)
			return nil, nil, err
		}
	} else {
		// Use transaction from the transaction manager
		cnx, _ = txManager.GetConnection()
	}
	defer cnx.Close()
	for _, update := range updates {
		var stmt *sql.Stmt
		stmt, err = cnx.PrepareContext(context.Background(), update.SQL)
		if err != nil {
			log.Printf("Error preparing update: %v", err)
			return nil, nil, err
		}

		if update.Values == nil {
			res, err := stmt.Exec()
			if err != nil {
				return nil, nil, err
			}

			rowCount, err := res.RowsAffected()
			if err != nil {
				return nil, nil, err
			}
			rowCounts = append(rowCounts, rowCount)

			newRowID, err := res.LastInsertId()
			if err != nil {
				newRowIDs = append(newRowIDs, 0) // Add 0 if no last ID is returned
			} else {
				newRowIDs = append(newRowIDs, newRowID)
			}
		} else {
			for _, values := range update.Values {
				res, err := stmt.Exec(values...)
				if err != nil {
					return nil, nil, err
				}

				rowCount, err := res.RowsAffected()
				if err != nil {
					return nil, nil, err
				}
				rowCounts = append(rowCounts, rowCount)

				newRowID, err := res.LastInsertId()
				if err != nil {
					newRowIDs = append(newRowIDs, 0)
				} else {
					newRowIDs = append(newRowIDs, newRowID)
				}
			}
		}
		stmt.Close()
	}

	return rowCounts, newRowIDs, nil
}

func (r *RDSPooledConnection) ExecuteFunctions(updateFunctions []func() error) error {
	r.txManagerPool.Register()

	defer func() {
		if rec := recover(); rec != nil {
			r.txManagerPool.Release(false)
			panic(rec)
		}
	}()

	for _, updateFunc := range updateFunctions {
		err := updateFunc()
		if err != nil {
			r.txManagerPool.Release(false)
			return err
		}
	}

	return r.txManagerPool.Release(true)
}
