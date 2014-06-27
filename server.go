package main

import (
	"database/sql"
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

type GeneralResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SetupDB() *sql.DB {
	// db, err := sql.Open("postgres", "user=root password=eaP7F1ZyCU6f40Ii host=172.17.42.1 port=49155 dbname=db sslmode=disable")
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

	m.Get("/wc_players", binding.Bind(PlayerNames{}), getWcPlayers)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
	m.Run()
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
