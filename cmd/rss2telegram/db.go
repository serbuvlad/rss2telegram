package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB

	stmtInsert *sql.Stmt
)

func initDB(dbPath string) {
	var err error

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		// Apparently this doesn't error out
		log.Fatal(fmt.Errorf("failed to open database %s: %w", dbPath, err))
	}

	stmtCreate, err := db.Prepare("CREATE TABLE IF NOT EXISTS read (id INTEGER PRIMARY KEY, feed TEXT, guid TEXT)")
	if err != nil {
		// Instead it errors out when we try to first do something;
		// this will give an inaccurate error if the statement above fails to compile
		log.Fatal(fmt.Errorf("failed to open database file: %w", err))

		//panic(fmt.Errorf("failed to compile create table statement: %w", err))
	}

	_, err = stmtCreate.Exec()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create table in database: %w", err))
	}

	stmtInsert, err = db.Prepare("INSERT INTO read (feed, guid) VALUES (?, ?)")
	if err != nil {
		panic(fmt.Errorf("failed to compile insert statmenet: %w", err))
	}
}

func addEntry(feed, guid string) {
	_, err := stmtInsert.Exec(feed, guid)

	if err != nil {
		log.Println(fmt.Errorf("failed to insert db values feed=%s guid=%s: %w", feed, guid, err))
	}
}

func haveEntry(feed, guid string) bool {
	have := false

	rows, err := db.Query("SELECT id FROM read WHERE feed = ? AND guid = ?", feed, guid)
	if err != nil {
		queryFail(feed, guid, err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		have = true
	}

	if err = rows.Err(); err != nil {
		queryFail(feed, guid, err)
		return false
	}

	return have
}

func queryFail(feed, guid string, err error) {
	log.Println(fmt.Errorf("query faliure feed=%s guid=%s: %w", feed, guid, err))

	log.Println("being conservative and refusing to send possible spam")
}
