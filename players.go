package main

import (
	"database/sql"
  "fmt"
	. "github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"unicode"
)

type PlayerStatsResponse struct {
	Status int          `json:"status"`
	Stats  []PlayerStat `json:"stats"`
}

type PlayerStat struct {
	Name     string `json:"name"`
	NameAlias   string `json:"nameAlias"`
	Club     string `json:"club"`
	Position string `json:"position"`
	Dob      string `json:"dob"`
	Height   string `json:"height"`
	Age      int64  `json:"age"`
	Weight   string `json:"weight"`
	National string `json:"national"`
	Image    string `json:"image"`

	Appearances string `json:"appearances"`
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
}

func GetPlayerStat(nameAlias string) (*PlayerStat, error) {
	// Load the HTML document (in real use, the type would be *goquery.Document)
	var statDoc, careerDoc, overviewDoc *Document
	var e error

	if statDoc, e = NewDocument("http://www.premierleague.com/en-gb/players/profile.statistics.html/" + nameAlias); e != nil {
		return nil, e
	}

	if overviewDoc, e = NewDocument("http://www.premierleague.com/en-gb/players/profile.overview.html/" + nameAlias); e != nil {
		return nil, e
	}

	if careerDoc, e = NewDocument("http://www.premierleague.com/en-gb/players/profile.career-history.html/" + nameAlias); e != nil {
		return nil, e
	}

	// general
	playerName := overviewDoc.Find(".hero-name .name").Eq(0).Text()
	club := overviewDoc.Find(".stats li").Eq(0).Find("p").Text()
	position := Captialize(strings.ToLower(overviewDoc.Find(".stats li").Eq(1).Find("p").Text()))
	dob := overviewDoc.Find(".contentTable .normal").Eq(0).Text()
	height := overviewDoc.Find(".contentTable .normal").Eq(1).Text()
	age, _ := strconv.ParseInt(overviewDoc.Find(".contentTable .normal").Eq(2).Text(), 10, 0)
	weight := overviewDoc.Find(".contentTable .normal").Eq(3).Text()
	national, _ := overviewDoc.Find(".contentTable .normal").Eq(5).Find("img").Attr("title")
	image, _ := overviewDoc.Find(".herosection .heroimg").Attr("src")

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

	playerStat := &PlayerStat{
		Name:        playerName,
		NameAlias:   nameAlias,
		Club:        club,
		Position:    position,
		Dob:         dob,
		Height:      height,
		Age:         age,
		Image:       image,
		Appearances: appearances,
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

	db, err := sql.Open("postgres", "user=plee dbname=fcards sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("INSERT INTO players (name, nameAlias, club, position, goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards, dob, height, age, weight, national, image, appearances) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)", playerName, nameAlias, club, position, goals, shots, penalties, assists, crosses, offsides, savesMade, ownGoals, cleanSheets, blocks, clearances, fouls, cards, dob, height, age, weight, national, image, appearances)
	if err != nil {
		fmt.Println("Insert error", err)
	}

	return playerStat, nil
}

func Captialize(str string) string {
	letters := []rune(str)
	letters[0] = unicode.ToUpper(letters[0])
	cappedStr := string(letters)
	return cappedStr
}
