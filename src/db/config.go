package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlConnPool *sql.DB
	TxManagerPool *TransactionManagerRegistry
	RdsConnPool   *RDSPooledConnection
)

func init() {
	var err error
	mysqlConnPool, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal(err)
	}

	TxManagerPool = NewTransactionManagerRegistry(mysqlConnPool)
	RdsConnPool = NewRDSPooledConnection(mysqlConnPool, TxManagerPool)

	fmt.Println("RDS Connection Pool and Transaction Manager initialized.")
}
