package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"unicode"
)

func main() {
	resp, err := http.Get("http://worldcup.kimonolabs.com/api/players?limit=800&apikey=0485cc45db1c8668a85c41d23ebff0b3")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	byt := []byte(string(body))
	var playerDat []map[string]interface{}

	if err := json.Unmarshal(byt, &playerDat); err != nil {
		panic(err)
	}

	resp, err = http.Get("http://worldcup.kimonolabs.com/api/clubs?limit=800&apikey=0485cc45db1c8668a85c41d23ebff0b3")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	byt = []byte(string(body))
	var clubDat []map[string]interface{}

	if err := json.Unmarshal(byt, &clubDat); err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "user=ins429 dbname=fcards sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(playerDat); i++ {
		for j := 0; j < len(clubDat); j++ {
			if playerDat[i]["clubId"] == clubDat[j]["id"] {
				playerDat[i]["club"] = clubDat[j]["name"]
				break
			}
		}
		InsertPlayerStat(db, playerDat[i])
		fmt.Println(playerDat[i])
	}

}

func InsertPlayerStat(db *sql.DB, player map[string]interface{}) {
	// general
	age := player["age"]
	foot := player["foot"]
	goals := player["goals"]
	club := player["club"]
	uid := player["id"]
	name := player["firstName"].(string) + " " + player["lastName"].(string)
	birthCountry := player["birthCountry"]
	birthCity := player["birthCity"]
	penaltyGoals := player["penaltyGoals"]
	birthDate := player["birthDate"]
	image := player["image"]
	height := player["heightCm"]
	assists := player["assists"]
	weight := player["weightKg"]
	national := player["nationality"]
	position := player["position"]
	ownGoals := player["ownGoals"]

	_, err := db.Exec("INSERT INTO wc_players (age, foot, goals, club, uid, name, birthCountry, birthCity, penaltyGoals, birthDate, image, height, assists, weight, national, position, ownGoals) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)", age, foot, goals, club, uid, name, birthCountry, birthCity, penaltyGoals, birthDate, image, height, assists, weight, national, position, ownGoals)
	if err != nil {
		fmt.Println("Insert error", err)
	}
}

func Capitalize(str string) string {
	letters := []rune(str)
	letters[0] = unicode.ToUpper(letters[0])
	cappedStr := string(letters)
	return cappedStr
}

func Uncapitalize(str string) string {
	letters := []rune(str)
	letters[0] = unicode.ToLower(letters[0])
	cappedStr := string(letters)
	return cappedStr
}
