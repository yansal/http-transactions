package main

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/yansal/http-transactions/handler"
	"github.com/yansal/http-transactions/manager"
	"github.com/yansal/sqldriver"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "sslmode=disable"
	}
	pqconnector, err := pq.NewConnector(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(&sqldriver.Connector{
		Connector: pqconnector,
		BeginTxFunc: func(opts driver.TxOptions, duration time.Duration, err error) {
			log.Printf("query=BEGIN opts=%+v duration=%s err=%v", opts, duration, err)
		},
		CommitFunc: func(duration time.Duration, err error) {
			log.Printf("query=COMMIT duration=%s err=%v", duration, err)
		},
		NextFunc: func(dest []driver.Value, duration time.Duration, err error) {
			log.Printf("dest=%+v duration=%s err=%v", dest, duration, err)
		},
		QueryContextFunc: func(query string, args []driver.NamedValue, duration time.Duration, err error) {
			log.Printf("query=%q args=%+v duration=%s err=%v", query, args, duration, err)
		},
		RollbackFunc: func(duration time.Duration, err error) {
			log.Printf("query=ROLLBACK duration=%s err=%v", duration, err)
		},
	})

	http.Handle("/users", handler.NewUser(manager.NewUser(db)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
