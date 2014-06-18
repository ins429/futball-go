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

	db, err := sql.Open("postgres", "user=plee dbname=fcards sslmode=disable")
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

  var idStr string
  rows, err := db.Query("SELECT uid FROM wc_players WHERE uid = $1", player["id"].(string))
  for rows.Next() {
    err = rows.Scan(&idStr)
  }

	if err != nil {
		fmt.Println("Insert error", err)
	}

  if idStr == "" {
    _, err := db.Exec("INSERT INTO wc_players (age, foot, goals, club, uid, name, birthCountry, birthCity, penaltyGoals, birthDate, image, height, assists, weight, national, position, ownGoals) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)", age, foot, goals, club, uid, name, birthCountry, birthCity, penaltyGoals, birthDate, image, height, assists, weight, national, position, ownGoals)
    if err != nil {
      fmt.Println("Insert error", err)
    }
    fmt.Println(player["firstName"].(string) + " does not exists")
  } else {
    _, err := db.Exec("UPDATE wc_players SET age = $1, foot = $2, goals = $3, club = $4, uid = $5, name = $6, birthCountry = $7, birthCity = $8, penaltyGoals = $9, birthDate = $10, image = $11, height = $12, assists = $13, weight = $14, national = $15, position = $16, ownGoals = $17 WHERE uid = $18", age, foot, goals, club, uid, name, birthCountry, birthCity, penaltyGoals, birthDate, image, height, assists, weight, national, position, ownGoals, uid)
    if err != nil {
      fmt.Println("Update error", err)
    }
    fmt.Println(player["firstName"].(string) + " exists")
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
