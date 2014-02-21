package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"os"
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

	m.Use(render.Renderer(render.Options{
		Delims: render.Delims{"{[{", "}]}"},
	}))

	m.Get("/", func(params martini.Params, r render.Render) {
		r.HTML(200, "index", nil)
	})

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

	m.Post("/images", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		defer file.Close()

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create("/tmp/file")
		if err != nil {
			fmt.Fprintf(w, "Failed to open the file for writing")
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		// the header contains useful info, like the original file name
		fmt.Fprintf(w, "File %s uploaded successfully.", header.Filename)
	})

	m.Run()
}
