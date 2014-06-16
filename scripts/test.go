package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=ins429 dbname=fcards sslmode=disable")
	if err != nil {
		panic(err)
	}
	var name string
	var names []string
	rows, err := db.Query("SELECT name FROM wc_players")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan: ", err)
		}
		names = append(names, name)
	}
	fmt.Println(names)
}
