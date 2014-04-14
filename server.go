package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Card struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GeneralResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CardResponse struct {
	Code  int    `json:"code"`
	Cards []Card `json:"cards"`
}

type UserForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	password string `json:"password"`
}

func main() {
	m := martini.Classic()

	db, err := sql.Open("postgres", "user=ins429 dbname=fcards sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	m.Use(render.Renderer(render.Options{
		Delims: render.Delims{"{[{", "}]}"},
	}))

	store := sessions.NewCookieStore([]byte("ins429"))
	m.Use(sessions.Sessions("peter", store))

	m.Get("/", func(params martini.Params, r render.Render) {
		r.HTML(200, "index", nil)
	})

	// query cards
	m.Get("/cards", func(params martini.Params, r render.Render, rw http.ResponseWriter, req *http.Request) {
		limit, _ := strconv.ParseInt(req.URL.Query().Get("limit"), 10, 0)
		skip, _ := strconv.ParseInt(req.URL.Query().Get("skip"), 10, 0)

		// set default skip to 10
		if limit < 10 {
			limit = 10
		}

		rows, err := db.Query("SELECT id, name FROM cards LIMIT $1 OFFSET $2", limit, skip)
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		// for the consistency on the response, put it in array
		cards := []Card{}
		for rows.Next() {
			var i Card
			err = rows.Scan(&i.Id, &i.Name)
			if err != nil {
				fmt.Println("Scan: ", err)
			}

			cards = append(cards, i)
		}

		// build response for cards
		res := &CardResponse{
			Code:  200,
			Cards: cards}

		r.JSON(200, res)
	})

	// get card by id
	m.Get("/cards/:id", func(params martini.Params, r render.Render) {
		// query by the card id
		rows, err := db.Query("SELECT id, name FROM cards WHERE id = $1", params["id"])
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		cards := []Card{}
		for rows.Next() {
			var i Card
			err = rows.Scan(&i.Id, &i.Name)
			if err != nil {
				fmt.Println("Scan: ", err)
			}

			cards = append(cards, i)
		}

		// build response for cards
		res := &CardResponse{
			Code:  200,
			Cards: cards}

		r.JSON(200, res)
	})

	// user login
	m.Post("/login", func(r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
		email, password := req.FormValue("email"), req.FormValue("password")

		rows, _ := db.Query("SELECT id, email, encrypted_password FROM users WHERE email = $1", email)
		var u User
		for rows.Next() {
			err = rows.Scan(&u.Id, &u.Email, &u.password)
			if err != nil {
				fmt.Println("Scan: ", err)
			}
		}

		pass := []byte(password)
		upass := []byte(u.password)

		if bcrypt.CompareHashAndPassword(upass, pass) == nil {
			s.Set("userId", u.Id)
			r.JSON(200, u)
		} else {
			r.JSON(400, &GeneralResponse{
				Code:    400,
				Message: "Failed to login!"})
		}
	})

	// user signup
	m.Post("/signup", func(r render.Render, req *http.Request) {
		email, password := req.FormValue("email"), req.FormValue("password")
		pass := []byte(password)
		p, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
		if err != nil {
			return
		}

		_, err = db.Exec("INSERT INTO users (email, encrypted_password) VALUES ($1, $2)", email, p)
		if err != nil {
			fmt.Println("Insert error", err)
			r.JSON(500, &GeneralResponse{
				Code:    500,
				Message: "Failed to signup!"})
			return
		}

		r.JSON(200, &GeneralResponse{
			Code:    200,
			Message: "Successfully sign up!"})
	})

	m.Post("/cards", func(rw http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		defer file.Close()

		if err != nil {
			fmt.Fprintln(rw, err)
			return
		}

		out, err := os.Create("/tmp/file")
		if err != nil {
			fmt.Fprintf(rw, "Failed to open the file for writing")
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(rw, err)
		}

		// the header contains useful info, like the original file name
		fmt.Fprintf(rw, "File %s uploaded successfully.", header.Filename)
	})

	m.Get("/showme", func(params martini.Params, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
		user := &User{}
		err := db.QueryRow("SELECT id, email from users where id=$1", s.Get("userId")).Scan(&user.Id, &user.Email)

		if err != nil {
			r.JSON(400, &GeneralResponse{
				Code:    400,
				Message: "Failed to look up!"})
			return
		}

		r.JSON(200, user)
	})

	m.Get("/players/:name", func(params martini.Params, r render.Render) {
		playerStat, _ := GetPlayerStat(params["name"])
		playerStats := []PlayerStat{}
		playerStats = append(playerStats, *playerStat)

		// build response for player stats
		res := &PlayerStatsResponse{
			Code:  200,
			Stats: playerStats}

		r.JSON(200, res)
	})

	http.ListenAndServe(":8080", m)
	m.Run()
}
