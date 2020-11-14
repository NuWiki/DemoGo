package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func insert(todo Todo) {
	postgresDb, err := sql.Open("postgres", "postgres://mydglzdw:3dtcdLkiIO93u2CQDE8SC2hIFP80WDjh@suleiman.db.elephantsql.com:5432/mydglzdw")
	if err != nil {
		log.Fatal("Connect to database error", err)

	}
	defer postgresDb.Close()

	row := postgresDb.QueryRow("INSERT INTO todos (title, status) values ($1,$2) RETURNING id", "buy bmw", "active")
	var id int

	err = row.Scan(&id)
	if err != nil {
		fmt.Println("cant scan id", err)
		return
	}

	fmt.Println("insert todo success : ", id)

}
