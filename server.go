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
	"log"
	"net/http"
	"os"
	"strconv"

	"strings"
	"unicode"

	// goquery
	. "github.com/PuerkitoBio/goquery"

	// for password
	"code.google.com/p/go.crypto/bcrypt"
	// "time"
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

type PlayerStatsResponse struct {
	Code  int          `json:"code"`
	Stats []PlayerStat `json:"stats"`
}

type PlayerStat struct {
	Club        string `json:"club"`
	Position    string `json:"position"`
	Goals       int64  `json:"goals"`
	Shots       int64  `json:"shots"`
	Penalties   int64  `json:"penalties"`
	Assists     int64  `json:"assists"`
	Crosses     int64  `json:"crosses"`
	Offsides    int64  `json:"offsides"`
	SavesMade   int64  `json:"savesMade"`
	OwnGoals    int64  `json:"ownGoals"`
	CleanSheets int64  `json:"cleanSheets"`
	Blocks      int64  `json:"blocks"`
	Clearances  int64  `json:"clearances"`
	Fouls       int64  `json:"fouls"`
	Cards       int64  `json:"cards"`
	Dob         string `json:"dob"`
	Height      string `json:"height"`
	Age         int64  `json:"age"`
	Weight      string `json:"weight"`
	National    string `json:"national"`
}

func GetPlayerStat(name string) (*PlayerStat, error) {
	// Load the HTML document (in real use, the type would be *goquery.Document)
	var statDoc, overviewDoc *Document
	var e error

	splitName := strings.Split(name, "-")
	fmt.Println(splitName)
	firstName := splitName[0]
	lastName := splitName[1]

	if statDoc, e = NewDocument("http://www.premierleague.com/en-gb/players/profile.statistics.html/" + firstName + "-" + lastName); e != nil {
		panic(e.Error())
	}

	if overviewDoc, e = NewDocument("http://www.premierleague.com/en-gb/players/profile.overview.html/" + firstName + "-" + lastName); e != nil {
		panic(e.Error())
	}

	// general
	club := overviewDoc.Find(".stats li").Eq(0).Find("p").Text()
	position := Captialize(strings.ToLower(overviewDoc.Find(".stats li").Eq(1).Find("p").Text()))
	dob := overviewDoc.Find(".contentTable .normal").Eq(0).Text()
	height := overviewDoc.Find(".contentTable .normal").Eq(1).Text()
	age, _ := strconv.ParseInt(overviewDoc.Find(".contentTable .normal").Eq(2).Text(), 10, 0)
	weight := overviewDoc.Find(".contentTable .normal").Eq(3).Text()
	national, _ := overviewDoc.Find(".contentTable .normal").Eq(5).Find("img").Attr("title")

	// attacking
	goals, _ := strconv.ParseInt(statDoc.Find("#clubsTabsAttacking li[name='goals'] .data").Text(), 10, 0)
	shots, _ := strconv.ParseInt(statDoc.Find("#clubsTabsAttacking li[name='shots'] .data").Text(), 10, 0)
	penalties, _ := strconv.ParseInt(statDoc.Find("#clubsTabsAttacking li[name='penalties'] .data").Text(), 10, 0)
	assists, _ := strconv.ParseInt(statDoc.Find("#clubsTabsAttacking li[name='assists'] .data").Text(), 10, 0)
	crosses, _ := strconv.ParseInt(statDoc.Find("#clubsTabsAttacking li[name='crosses'] .data").Text(), 10, 0)
	offsides, _ := strconv.ParseInt(statDoc.Find("#clubsTabsAttacking li[name='offsides'] .data").Text(), 10, 0)

	// defending
	savesMade, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDefending li[name='savesMade'] .data").Text(), 10, 0)
	ownGoals, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDefending li[name='ownGoals'] .data").Text(), 10, 0)
	cleanSheets, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDefending li[name='cleanSheets'] .data").Text(), 10, 0)
	blocks, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDefending li[name='blocks'] .data").Text(), 10, 0)
	clearances, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDefending li[name='clearances'] .data").Text(), 10, 0)

	// disciplinary
	fouls, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDisciplinary li[name='fouls'] .data").Text(), 10, 0)
	cards, _ := strconv.ParseInt(statDoc.Find("#clubsTabsDisciplinary li[name='cards'] .data").Text(), 10, 0)

	fmt.Println(goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards)

	playerStat := &PlayerStat{
		Club:        club,
		Position:    position,
		Dob:         dob,
		Height:      height,
		Age:         age,
		Weight:      weight,
		National:    national,
		Goals:       goals,
		Shots:       shots,
		Penalties:   penalties,
		Assists:     assists,
		Crosses:     crosses,
		Offsides:    offsides,
		SavesMade:   savesMade,
		OwnGoals:    ownGoals,
		CleanSheets: cleanSheets,
		Blocks:      blocks,
		Clearances:  clearances,
		Fouls:       fouls,
		Cards:       cards}

	return playerStat, nil
}

func Captialize(str string) string {
	letters := []rune(str)
	letters[0] = unicode.ToUpper(letters[0])
	cappedStr := string(letters)
	return cappedStr
}
