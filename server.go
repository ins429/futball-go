package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

type Card struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GeneralResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CardResponse struct {
	Status int    `json:"status"`
	Cards  []Card `json:"cards"`
}

type AddCardForm struct {
	Players []string `json:"players"`
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
			Status: 200,
			Cards:  cards}

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
			Status: 200,
			Cards:  cards}

		r.JSON(200, res)
	})

	// user login
	m.Post("/login", binding.Bind(UserForm{}), func(r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, userForm UserForm) {
		rows, _ := db.Query("SELECT id, username, encrypted_password FROM users WHERE username = $1", userForm.Username)
		var u User
		for rows.Next() {
			err = rows.Scan(&u.Id, &u.Username, &u.Password)
			if err != nil {
				fmt.Println("Scan: ", err)

				r.JSON(400, &GeneralResponse{
					Status:  400,
					Message: "Failed to login!"})
				return
			}
		}

		pass := []byte(userForm.Password)
		upass := []byte(u.Password)

		if bcrypt.CompareHashAndPassword(upass, pass) == nil {
			s.Set("userId", u.Id)
			r.JSON(200, &UserResponse{
				Status: 200,
				User:   u})
		} else {
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to login!"})
		}
	})

	m.Delete("/logout", func(r render.Render, s sessions.Session) {
		s.Delete("userId")

		r.JSON(200, &GeneralResponse{
			Status:  200,
			Message: "Successfully logged out!"})
	})

	// user signup
	m.Post("/signup", func(r render.Render, req *http.Request) {
		username, password := req.FormValue("username"), req.FormValue("password")
		pass := []byte(password)
		p, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
		if err != nil {
			return
		}

		_, err = db.Exec("INSERT INTO users (username, encrypted_password) VALUES ($1, $2)", username, p)
		if err != nil {
			fmt.Println("Insert error", err)
			r.JSON(500, &GeneralResponse{
				Status:  500,
				Message: "Failed to signup!"})
			return
		}

		r.JSON(200, &GeneralResponse{
			Status:  200,
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
		var players string
		err := db.QueryRow("SELECT id, username, firstname, lastname, players from users where id=$1", s.Get("userId")).Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &players)
		if err != nil {
			fmt.Println(err)
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to look up!"})
			return
		}
		newArray := []interface{}{"{\"name\":\"luis-suarez\"}", "{\"name\":\"leighton-baines\"}"}
		fmt.Println(newArray)
		test := interface{}(players)
		fmt.Println([]interface{}(test))
		playersByt := []byte(players)
		var dat []interface{}
		if err := json.Unmarshal(playersByt, &dat); err != nil {
			panic(err)
		}
		fmt.Println(reflect.TypeOf(dat))
		r.JSON(200, &UserResponse{
			Status: 200,
			User:   *user})
	})

	m.Get("/players/:name", func(params martini.Params, r render.Render) {
		playerStat, _ := GetPlayerStat(params["name"])
		playerStats := []PlayerStat{}
		playerStats = append(playerStats, *playerStat)

		// build response for player stats
		res := &PlayerStatsResponse{
			Status: 200,
			Stats:  playerStats}

		r.JSON(200, res)
	})

	m.Put("/add_card", binding.Bind(AddCardForm{}), func(r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, addCardForm AddCardForm) {
		user := &User{}
		err := db.QueryRow("SELECT players from users where id=$1", s.Get("userId")).Scan(&user.Players)

		fmt.Println(user.Players)
		for i := 0; i < len(addCardForm.Players); i++ {
			var player map[string]interface{}
			err = json.Unmarshal([]byte(addCardForm.Players[i]), &player)
		}

		fmt.Println(addCardForm.Players)
		_, err = db.Exec("UPDATE users SET players = $1 WHERE id = $2", addCardForm.Players, s.Get("userId"))
		if err != nil {
			fmt.Println("Insert error", err)
			r.JSON(500, &GeneralResponse{
				Status:  500,
				Message: "Failed to signup!"})
			return
		}

		r.JSON(200, &GeneralResponse{
			Status:  200,
			Message: "...!"})
	})

	m.Post("/fb_signup", binding.Bind(FbForm{}), func(r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, fbForm FbForm) {
		fbUser, err := FbGetMe(fbForm.Token)

		if err != nil {
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to signup!"})
			return
		}

		// check if the user exists
		rows, err := db.Query("SELECT id, username, lastname, firstname FROM users WHERE fb_id = $1", fbUser.Id)

		if err != nil {
			fmt.Println(err)
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to signup!"})
			return
		}

		if rows.Next() {
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "User already exists!"})
			return
		}

		_, err = db.Exec("INSERT INTO users (fb_id, username, firstname, lastname) VALUES ($1, $2, $3, $4)", fbUser.Id, fbUser.Username, fbUser.FirstName, fbUser.LastName)
		if err != nil {
			fmt.Println("Insert error", err)
			r.JSON(500, &GeneralResponse{
				Status:  500,
				Message: "Failed to signup!"})
			return
		}

		rows, err = db.Query("SELECT id, username, lastname, firstname FROM users WHERE fb_id = $1", fbUser.Id)

		if err != nil {
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to signup!"})
			return
		}

		var u User
		for rows.Next() {
			err = rows.Scan(&u.Id, &u.Username, &u.LastName, &u.FirstName)
			if err != nil {
				fmt.Println("Scan: ", err)

				r.JSON(400, &GeneralResponse{
					Status:  400,
					Message: "Failed to login!"})
				return
			}
		}

		s.Set("userId", u.Id)
		r.JSON(200, &UserResponse{
			Status: 200,
			User:   u})
	})

	m.Post("/fb_login", binding.Bind(FbForm{}), func(r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, fbForm FbForm) {
		fbUser, err := FbGetMe(fbForm.Token)

		if err != nil {
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to login!"})
			return
		}

		rows, err := db.Query("SELECT id, username, lastname, firstname FROM users WHERE fb_id = $1", fbUser.Id)
		if err != nil {
			fmt.Println(err)
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to login!"})
			return
		}

		var u User
		for rows.Next() {
			err = rows.Scan(&u.Id, &u.Username, &u.LastName, &u.FirstName)
			if err != nil {
				fmt.Println("Scan: ", err)

				r.JSON(400, &GeneralResponse{
					Status:  400,
					Message: "Failed to login!"})
				return
			}
		}

		s.Set("userId", u.Id)
		r.JSON(200, &UserResponse{
			Status: 200,
			User:   u})
	})

	http.ListenAndServe(":8080", m)
	m.Run()
}

func FbGetMe(token string) (FbUser, error) {
	fmt.Println("Getting me")
	response, err := getUncachedResponse("https://graph.facebook.com/me?access_token=" + token)

	if err == nil {
		responseBody := readHttpBody(response)

		if responseBody != "" {
			var fbUser FbUser
			err = json.Unmarshal([]byte(responseBody), &fbUser)

			if err == nil {
				return fbUser, nil
			}
		}
		return FbUser{}, err
	}

	return FbUser{}, err
}

func getUncachedResponse(uri string) (*http.Response, error) {
	fmt.Println("Uncached response GET")
	request, err := http.NewRequest("GET", uri, nil)

	if err == nil {
		request.Header.Add("Cache-Control", "no-cache")

		client := new(http.Client)

		return client.Do(request)
	}

	if err != nil {
	}
	return nil, err
}

func readHttpBody(response *http.Response) string {
	fmt.Println("Reading body")

	bodyBuffer := make([]byte, 1000)
	var str string

	count, err := response.Body.Read(bodyBuffer)

	for ; count > 0; count, err = response.Body.Read(bodyBuffer) {
		if err != nil {

		}

		str += string(bodyBuffer[:count])
	}

	return str
}
