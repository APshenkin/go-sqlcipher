package main

import (
	"database/sql"
	"fmt"
	"log"

	sqlite3 "github.com/mutecomm/go-sqlcipher/v4"
)

func main() {
	sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("series", &seriesModule{})
		},
	})
	db, err := sql.Open("sqlite3_with_extensions", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM series")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var value int
		rows.Scan(&value)
		fmt.Printf("value: %d\n", value)
	}
}
