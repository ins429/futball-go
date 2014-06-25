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
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AddCardForm struct {
	Name string `form:"name" json:"name"`
}

func SetupDB() *sql.DB {
	dokkuDB := os.Getenv("DATABASE_URL")
	fmt.Println("here")
	fmt.Println(dokkuDB)
	db, err := sql.Open("postgres", "user=root password=eaP7F1ZyCU6f40Ii host=172.17.42.1 port=49155 dbname=db sslmode=disable")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return db
}

func main() {
	fmt.Println("starting futbol-cards")
	m := martini.Classic()
	m.Map(SetupDB())

	m.Use(render.Renderer(render.Options{
		Delims: render.Delims{"{[{", "}]}"},
	}))

	store := sessions.NewCookieStore([]byte("ins429"))
	m.Use(sessions.Sessions("peter", store))

	m.Get("/", func(params martini.Params, r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Post("/fb_login", binding.Bind(FbForm{}), fbUserLogin)
	m.Post("/fb_signup", binding.Bind(FbForm{}), fbUserSignUp)
	m.Post("/login", binding.Bind(UserForm{}), userLogin)
	m.Delete("/logout", userLogout)
	m.Post("/signup", userSignUp)
	m.Get("/showme", userShowMe)

	m.Get("/wc_players", binding.Bind(PlayerNames{}), getWcPlayers)
	m.Get("/players", binding.Bind(PlayerNames{}), getPlayers)

	m.Put("/add_card", binding.Bind(AddCardForm{}), addCard)
	m.Delete("/remove_card", binding.Bind(AddCardForm{}), removeCard)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
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

func fbUserLogin(db *sql.DB, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, fbForm FbForm) {
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
}

func userLogin(db *sql.DB, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, userForm UserForm) {
	rows, _ := db.Query("SELECT id, username, encrypted_password FROM users WHERE username = $1", userForm.Username)
	var u User
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Username, &u.Password)
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
}

func userLogout(db *sql.DB, r render.Render, s sessions.Session) {
	s.Delete("userId")

	r.JSON(200, &GeneralResponse{
		Status:  200,
		Message: "Successfully logged out!"})
}

func userSignUp(db *sql.DB, r render.Render, req *http.Request) {
	username, password := req.FormValue("username"), req.FormValue("password")
	pass := []byte(password)
	p, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO users (username, encrypted_password) VALUES ($1, $2)", username, p)
	if err != nil {
		fmt.Println("Insert error", err)
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to signup!"})
		return
	}

	r.JSON(200, &GeneralResponse{
		Status:  200,
		Message: "Successfully sign up!"})
}

func userShowMe(db *sql.DB, params martini.Params, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	user := &User{}
	var playersRaw string
	err := db.QueryRow("SELECT id, username, firstname, lastname, array_to_json(players) from users where id=$1", s.Get("userId")).Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &playersRaw)
	if err != nil {
		fmt.Println(err)
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to look up!"})
		return
	}

	playersByt := []byte(playersRaw)
	var dat []interface{}
	if err := json.Unmarshal(playersByt, &dat); err != nil {
		panic(err)
	}
	user.Players = dat

	r.JSON(200, &UserResponse{
		Status: 200,
		User:   *user})
}

func getWcPlayers(db *sql.DB, params martini.Params, r render.Render, playerNames PlayerNames) {
	if len(playerNames.Names) == 0 {
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Please pass in player names!"})
		return
	}
	dollars := ""
	playerStats := []WCPlayerStat{}

	for i := 0; i < len(playerNames.Names); i++ {
		dollars += "$" + strconv.Itoa(i+1)
		if i < len(playerNames.Names)-1 {
			dollars += ","
		}
	}

	args := make([]interface{}, len(playerNames.Names))
	for i, s := range playerNames.Names {
		args[i] = s
	}
	fmt.Println(playerNames.Names)
	var name, foot, birthDate, birthCountry, birthCity, national, position, club, image []byte
	var age, height, weight, goals, assists, penaltyGoals, ownGoals float64
	rows, err := db.Query("SELECT name, age, foot, birthDate, birthCountry, birthCity, national, height, weight, position, club, goals, assists, penaltyGoals, ownGoals, image FROM wc_players WHERE name IN ("+dollars+")", args...)

	if err != nil {
		fmt.Println("Query: ", err)
	}

	for rows.Next() {
		err = rows.Scan(&name, &age, &foot, &birthDate, &birthCountry, &birthCity, &national, &height, &weight, &position, &club, &goals, &assists, &penaltyGoals, &ownGoals, &image)
		if err != nil {
			fmt.Println("Scan: ", err)

			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to get players!"})
			return
		}

		p := &WCPlayerStat{
			Name:         string(name),
			Age:          age,
			Foot:         string(foot),
			BirthDate:    string(birthDate),
			BirthCountry: string(birthCountry),
			BirthCity:    string(birthCity),
			National:     string(national),
			Height:       height,
			Weight:       weight,
			Position:     string(position),
			Club:         string(club),
			Goals:        goals,
			Assists:      assists,
			PenaltyGoals: penaltyGoals,
			OwnGoals:     ownGoals,
			Image:        string(image)}
		playerStats = append(playerStats, *p)
	}

	// build response for player stats
	res := &WCPlayerStatsResponse{
		Status: 200,
		Stats:  playerStats}

	r.JSON(200, res)
}

func getPlayers(db *sql.DB, params martini.Params, r render.Render, playerNames PlayerNames) {
	if len(playerNames.Names) == 0 {
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Please pass in player names!"})
		return
	}
	dollars := ""
	playerStats := []PlayerStat{}

	for i := 0; i < len(playerNames.Names); i++ {
		dollars += "$" + strconv.Itoa(i+1)
		if i < len(playerNames.Names)-1 {
			dollars += ","
		}
	}

	args := make([]interface{}, len(playerNames.Names))
	for i, s := range playerNames.Names {
		args[i] = s
	}

	rows, err := db.Query("SELECT name, nameAlias, club, position, dob, height, age, weight, national, image, appearances, goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards FROM players WHERE nameAlias IN ("+dollars+")", args...)

	if err != nil {
		fmt.Println("Query: ", err)
	}

	var p PlayerStat
	for rows.Next() {
		err = rows.Scan(&p.Name, &p.NameAlias, &p.Club, &p.Position, &p.Dob, &p.Height, &p.Age, &p.Weight, &p.National, &p.Image, &p.Appearances, &p.Goals, &p.Shots, &p.Penalties, &p.Assists, &p.Crosses, &p.Offsides, &p.SavesMade, &p.OwnGoals, &p.CleanSheets, &p.Blocks, &p.Clearances, &p.Fouls, &p.Cards)
		if err != nil {
			fmt.Println("Scan: ", err)

			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to get players!"})
			return
		}

		playerStats = append(playerStats, p)
	}

	// build response for player stats
	res := &PlayerStatsResponse{
		Status: 200,
		Stats:  playerStats}

	r.JSON(200, res)
}

func addCard(db *sql.DB, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, addCardForm AddCardForm) {
	// check user session
	if s.Get("userId") == nil || s.Get("userId") == "" {
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Please login first!"})
		return
	}

	var playersRaw string
	err := db.QueryRow("SELECT array_to_json(players) from users where id=$1", s.Get("userId")).Scan(&playersRaw)
	if err != nil {
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to add a card!"})
		return
	}

	playersByt := []byte(playersRaw)
	var userPlayers []map[string]interface{}
	if err := json.Unmarshal(playersByt, &userPlayers); err != nil {
		fmt.Println(err)
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to add a card!"})
		return
	}

	for i := 0; i < len(userPlayers); i++ {
		// if the name already exists in user's card list
		if userPlayers[i]["name"] == addCardForm.Name {
			r.JSON(400, &GeneralResponse{
				Status:  400,
				Message: "Failed to add a card, " + addCardForm.Name + " already exists!"})
			return
		}
	}

	_, err = db.Exec("UPDATE users SET players = array_append(players, $1) WHERE id = $2", "{\"name\":\""+addCardForm.Name+"\"}", s.Get("userId"))
	if err != nil {
		fmt.Println("Update error", err)
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to add a card!"})
		return
	}

	r.JSON(200, &GeneralResponse{
		Status:  200,
		Message: "...!"})
}

func removeCard(db *sql.DB, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, addCardForm AddCardForm) {
	var playersRaw string
	err := db.QueryRow("SELECT array_to_json(players) from users where id=$1", s.Get("userId")).Scan(&playersRaw)
	if err != nil {
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to add a card!"})
		return
	}

	playersByt := []byte(playersRaw)
	var userPlayers []map[string]interface{}

	if err := json.Unmarshal(playersByt, &userPlayers); err != nil {
		fmt.Println(err)
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to add a card!"})
		return
	}

	playerFound := false
	dollars := ""
	for i := 0; i < len(userPlayers); i++ {
		// if the name already exists in user's card list
		if userPlayers[i]["name"] == addCardForm.Name {
			copy(userPlayers[i:], userPlayers[i+1:])
			userPlayers[len(userPlayers)-1] = nil // or the zero value of T
			userPlayers = userPlayers[:len(userPlayers)-1]
			playerFound = true
		}
	}

	for i := 0; i < len(userPlayers); i++ {
		dollars += "$" + strconv.Itoa(i+1)
		if i < len(userPlayers)-1 {
			dollars += ","
		}
	}

	if !playerFound || err != nil {
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to remove a card!"})
		return
	}

	args := make([]interface{}, len(userPlayers))
	for i, s := range userPlayers {
		userPlayersJson, _ := json.Marshal(s)
		args[i] = string(userPlayersJson)
	}
	args = append(args, s.Get("userId"))

	_, err = db.Exec("UPDATE users SET players = CAST(ARRAY["+dollars+"] as json[]) WHERE id = $"+strconv.Itoa(len(userPlayers)+1), args...)
	if err != nil {
		fmt.Println("Update error", err)
		r.JSON(400, &GeneralResponse{
			Status:  400,
			Message: "Failed to add a card!"})
		return
	}

	r.JSON(200, &GeneralResponse{
		Status:  200,
		Message: "Removed!"})
}

func fbUserSignUp(db *sql.DB, r render.Render, rw http.ResponseWriter, req *http.Request, s sessions.Session, fbForm FbForm) {
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
		r.JSON(400, &GeneralResponse{
			Status:  400,
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
}
