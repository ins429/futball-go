package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/lib/pq"
)

type Image struct {
	Name string
	Url  string
}

func main() {
	m := martini.Classic()
	// m.Get("/", func() string {
	// 	return "Hello world!"
	// })

	db, err := sql.Open("postgres", "user=ins429 dbname=futball_gifs_development sslmode=disable")
	if err != nil {
		fmt.Print(err)
	}

	m.Use(render.Renderer())
	m.Get("/images/:id", func(params martini.Params, r render.Render) {
		rows, err := db.Query("SELECT * FROM images WHERE id = $1", params["id"])
		if err != nil {
			fmt.Print(err)
		}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				fmt.Print(err)
			}
			fmt.Printf(name)
		}
		r.JSON(200, Image{"hello", "world!"})
	})

	m.Run()
}
