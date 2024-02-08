package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	"log"
	"os"
)

type Tag struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func (t *Tag) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), t)
}

func (t *Tag) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func main() {
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE foo (tag JSONB)`)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO foo(tag) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(`{"name": "mattn", "country": "japan"}`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(`{"name": "michael", "country": "usa"}`)
	if err != nil {
		log.Fatal(err)
	}

	var country string
	err = db.QueryRow("SELECT tag->>'country' FROM foo WHERE tag->>'name' = 'mattn'").Scan(&country)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(country)

	var tag Tag
	err = db.QueryRow("SELECT tag FROM foo WHERE tag->>'name' = 'mattn'").Scan(&tag)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tag.Name)

	tag.Country = "日本"
	_, err = db.Exec(`UPDATE foo SET tag = ? WHERE tag->>'name' == 'mattn'`, &tag)
	if err != nil {
		log.Fatal(err)
	}

	err = db.QueryRow("SELECT tag->>'country' FROM foo WHERE tag->>'name' = 'mattn'").Scan(&country)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(country)
}
