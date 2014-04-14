package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"strings"
	"unicode"

	// goquery
	. "github.com/PuerkitoBio/goquery"
)

func main() {
	for i := 1; i <= 241; i++ {
		resp, err := http.Get("http://www.premierleague.com/ajax/player/index/A_TO_Z/null/null/ALL/null/null/null/ALL/null/null/20/4/2/2/" + strconv.Itoa(i) + "/null.json")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		byt := []byte(string(body))

		var dat map[string]interface{}

		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		playersList := dat["playerIndexSection"].(map[string]interface{})["index"].(map[string]interface{})["resultsList"].([]interface{})
		for j := 0; j < len(playersList); j++ {
			if playersList[j].(map[string]interface{})["activeInPremierLeague"] == true {
				name := playersList[j].(map[string]interface{})["cmsAlias"].([]interface{})[0].(string)

				GetPlayerStat(name)
			}
		}

		if err != nil {
			log.Fatal(err)
		}

	}
}

func GetPlayerStat(name string) {
	// Load the HTML document (in real use, the type would be *goquery.Document)
	var statDoc, careerDoc, overviewDoc *Document
	var err error

	if statDoc, err = NewDocument("http://www.premierleague.com/en-gb/players/profile.statistics.html/" + name); err != nil {
		panic(err.Error())
	}

	if careerDoc, err = NewDocument("http://www.premierleague.com/en-gb/players/profile.career-history.html/" + name); err != nil {
		panic(err.Error())
	}

	if overviewDoc, err = NewDocument("http://www.premierleague.com/en-gb/players/profile.overview.html/" + name); err != nil {
		panic(err.Error())
	}

	// general
	playerName := overviewDoc.Find(".hero-name .name").Eq(0).Text()

	if playerName == "" {
		fmt.Println("error getting player:" + name)
		return
	}
	club := overviewDoc.Find(".stats li").Eq(0).Find("p").Text()
	position := Capitalize(strings.ToLower(overviewDoc.Find(".stats li").Eq(1).Find("p").Text()))
	dob := overviewDoc.Find(".contentTable .normal").Eq(0).Text()
	height := overviewDoc.Find(".contentTable .normal").Eq(1).Text()
	age, _ := strconv.ParseInt(overviewDoc.Find(".contentTable .normal").Eq(2).Text(), 10, 0)
	weight := overviewDoc.Find(".contentTable .normal").Eq(3).Text()
	national, _ := overviewDoc.Find(".contentTable .normal").Eq(5).Find("img").Attr("title")

	appearances := strings.Replace(careerDoc.Find(".contentTable.stats").Eq(0).Find("tr:nth-child(2) td:nth-child(2)").Text(), "\t", "", -1)
	appearances = strings.Replace(appearances, "\n", "", -1)
	appearances = strings.Replace(appearances, " ", "", -1)

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

	fmt.Println(playerName, goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards, appearances)

	db, err := sql.Open("postgres", "user=ins429 dbname=fcards sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO players (name, nameAlias, club, position, appearances, goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards, dob, height, age, weight, national) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)", playerName, name, club, position, appearances, goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards, dob, height, age, weight, national)
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
