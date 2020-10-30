package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// COVidStatPart contains region statistic
type COVidStatPart struct {
	sick   float64
	healed float64
	died   float64
}

// COVidStat contain region name and region statistic
type COVidStat struct {
	region string
	part   COVidStatPart
}

func processCoronaQuery(message string) string {

	// http request section
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	covidURL := "http://covid19.bvn13.com/stats/last"

	log.Print("sending")
	req, err := http.NewRequest(http.MethodGet, covidURL, nil)
	_check(err)

	log.Print("wait")
	res, err := spaceClient.Do(req)
	_check(err)

	log.Print("read")
	body, err := ioutil.ReadAll(res.Body)
	_check(err)

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	list := dat["stats"].([]interface{})

	log.Println(SystemConfig.dbPath)
	// open base
	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	_check(err)
	defer db.Close()

	log.Println("Start find state")
	if checkKey(message, db) == true {
		message = strings.ToLower(getValueState(message, db))
	}

	for i := 0; i < len(list); i++ {
		currentList := list[i].(map[string]interface{})

		if strings.ToLower(currentList["region"].(string)) == message {
			log.Print("Region " + strconv.Itoa(i) + " : ")
			currentState := COVidStatPart{currentList["sick"].(float64), currentList["healed"].(float64), currentList["died"].(float64)}
			log.Println(currentList["region"].(string))
			log.Println(currentState)

			prevList := currentList["previous"].(map[string]interface{})
			prevState := COVidStatPart{prevList["sick"].(float64), prevList["healed"].(float64), prevList["died"].(float64)}
			log.Println(prevState)
			log.Println("Increment " + strconv.FormatFloat(currentState.sick-prevState.sick, 'f', 0, 64))

			return "В злоебучке " + message + "\n" +
				strconv.FormatFloat(currentState.died-prevState.died, 'f', 0, 64) + " \tкончилось\n" +
				strconv.FormatFloat(currentState.healed-prevState.healed, 'f', 0, 64) + " похорошело\n" +
				strconv.FormatFloat(currentState.sick-prevState.sick, 'f', 0, 64) + " подхватило"
		}
	}

	//	sublist := list[0].(map[string]interface{})

	return "ты хуйло."
}
