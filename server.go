package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	// for password
	"code.google.com/p/go.crypto/bcrypt"
	// "time"
	"code.google.com/p/go.net/html"
)

type Card struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
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

	m.Get("/", func(params martini.Params, r render.Render) {
		r.HTML(200, "index", nil)
	})

	// query cards
	m.Get("/cards", func(params martini.Params, r render.Render, w http.ResponseWriter, req *http.Request) {
		ExampleParse()
		checkSession(req, w)

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
		ExampleParse()
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
	m.Post("/login", binding.Form(UserForm{}), func(userForm UserForm, r render.Render, w http.ResponseWriter, req *http.Request) {

		rows, _ := db.Query("SELECT id, email, encrypted_password FROM users WHERE email = $1", userForm.Email)
		var u User
		for rows.Next() {
			err = rows.Scan(&u.Id, &u.Email, &u.password)
			if err != nil {
				fmt.Println("Scan: ", err)
			}
		}

		pass := []byte(userForm.Password)
		upass := []byte(u.password)

		fmt.Println(bcrypt.CompareHashAndPassword(upass, pass))

		if bcrypt.CompareHashAndPassword(upass, pass) == nil {
			saveSession(req, w, u.Id)
			r.JSON(200, u)
		} else {
			r.JSON(400, &ErrorResponse{
				Code:    400,
				Message: "Failed to login!"})
		}
	})

	// user signup
	m.Post("/signup", binding.Form(UserForm{}), func(userForm UserForm, r render.Render) {
		pass := []byte(userForm.Password)
		p, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(pass)
			fmt.Println(userForm.Password)
			fmt.Println(err)
			fmt.Println("umm.. error on GenerateFromPassword")
			return
		}

		_, err = db.Exec("INSERT INTO users (email, encrypted_password) VALUES ($1, $2)", userForm.Email, p)
		if err != nil {
			fmt.Println("Insert error", err)
			r.JSON(500, userForm)
			return
		}

		r.JSON(200, userForm)
	})

	m.Post("/cards", func(w http.ResponseWriter, r *http.Request) {
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

	m.Get("/showme", func(params martini.Params, r render.Render, w http.ResponseWriter, req *http.Request) {
		checkSession(req, w)

		r.JSON(200, UserForm{})
	})

	http.ListenAndServe(":8080", m)
	m.Run()
}

func checkSession(req *http.Request, rsp http.ResponseWriter) {
	store := NewPGStore("user=ins429 dbname=fcards sslmode=disable", []byte("something-very-secret"))

	defer store.Close()

	// Get a session.
	session, err := store.Get(req, "session-key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(session)

	fmt.Println(session.Values["user_id"])
	// Add a value.
	// session.Values["user_id"] = "1"

	// Save.
	if err = sessions.Save(req, rsp); err != nil {
		fmt.Println(err)
		log.Fatalf("Error saving session: %v", err)
	}

	// Delete session.
	// session.Options.MaxAge = -1
	// if err = sessions.Save(req, rsp); err != nil {
	// 	fmt.Println(err)
	// 	log.Fatalf("Error saving session: %v", err)
	// }
}

func isUserLogged(req *http.Request, rsp http.ResponseWriter) bool {
	store := NewPGStore("user=ins429 dbname=fcards sslmode=disable", []byte("something-very-secret"))

	defer store.Close()

	// Get a session.
	session, err := store.Get(req, "session-key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(session)

	fmt.Println(session.Values["user_id"])
	// Add a value.
	// session.Values["user_id"] = "1"

	// Save.
	if err = sessions.Save(req, rsp); err != nil {
		fmt.Println(err)
		log.Fatalf("Error saving session: %v", err)
	}

	return session.Values["user_id"] != nil
}

func saveSession(req *http.Request, rsp http.ResponseWriter, userId int64) {
	store := NewPGStore("user=ins429 dbname=fcards sslmode=disable", []byte("something-very-secret"))

	defer store.Close()

	// Get a session.
	session, err := store.Get(req, "session-key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(session)

	fmt.Println(session.Values["user_id"])
	// Add a value.
	session.Values["user_id"] = userId

	// Save.
	if err = sessions.Save(req, rsp); err != nil {
		fmt.Println(err)
		log.Fatalf("Error saving session: %v", err)
	}
}

func ExampleParse() {
	resp, err := http.Get("http://www.google.com/")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	s := string(body)
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	// Output:
	// foo
	// /bar/baz
}
