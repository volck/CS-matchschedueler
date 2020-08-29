package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MatchConfig struct {
	Matchid    string `json:"matchid"`
	NumMaps    int    `json:"num_maps"`
	Spectators struct {
		Players []string `json:"players"`
	} `json:"spectators"`
	Maplist []string `json:"maplist"`
	Team1   struct {
		Name    string `json:"name"`
		Tag     string `json:"tag"`
		Flag    string `json:"flag"`
		Logo    string `json:"logo"`
		Players struct {
			STEAM0152245092 string `json:"STEAM_0:1:52245092"`
			STEAM11         string `json:"STEAM_1:1:....."`
		} `json:"players"`
	} `json:"team1"`
	Team2 struct {
		Name    string   `json:"name"`
		Tag     string   `json:"tag"`
		Flag    string   `json:"flag"`
		Logo    string   `json:"logo"`
		Players []string `json:"players"`
	} `json:"team2"`
	Cvars struct {
		Hostname string `json:"hostname"`
	} `json:"cvars"`
	Matchtime string `json:"matchtime"`
}

var matchConfigs []MatchConfig

func makeMatch(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var match MatchConfig
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&match)
	if err != nil {
		json.NewEncoder(rw).Encode(err)
		fmt.Println(err.Error())
	} else {
		match.Cvars.Hostname = fmt.Sprintf("matchserver: %s vs %s", match.Team1.Name, match.Team2.Name)
		json.NewEncoder(rw).Encode(fmt.Sprintf("%s vs %s going to be played on %s schedueled for %s", match.Team1.Name, match.Team2.Name, match.Cvars.Hostname, match.Matchtime))
		fmt.Sprintf("%s vs %s going to be played on %s schedueled for %s", match.Team1.Name, match.Team2.Name, match.Cvars.Hostname, match.Matchtime)

		matchConfigs = append(matchConfigs, match)
	}

}

func unschedueledMatches(rw http.ResponseWriter, req *http.Request) {
	if len(matchConfigs) == 0 {
		json.NewEncoder(rw).Encode(false)

	} else {
		json.NewEncoder(rw).Encode(true)

	}
}

func scheduleMatches(rw http.ResponseWriter, req *http.Request) {
	layout := "02-01-2006 15:04:05"
	for match := range matchConfigs {
		parsedTime, err := time.Parse(layout, matchConfigs[match].Matchtime)
		if err != nil {
			fmt.Println(err)

		}
		layout := "02-01-2006 15:04:05"
		testTime := matchConfigs[match].Matchtime
		parsedTime, err = time.Parse(layout, testTime)
		loc, _ := time.LoadLocation("Europe/Oslo")
		timestamps := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), loc)
		if err != nil {
			fmt.Println(err)
		}

		now := time.Now()
		when := time.Until(timestamps)

		if when < time.Hour {
			if when < 0 {
				json.NewEncoder(rw).Encode(fmt.Sprintf("[%s * ] Match which was schedueled for %v is stale!", now.In(loc), testTime))


			} else {
				json.NewEncoder(rw).Encode(fmt.Sprintf("[%s * ]  match(%s) needs to be scheduled within the hour: %v till gametime", now.In(loc), testTime, time.Until(timestamps)))

			}
		} else if when > time.Hour {
			json.NewEncoder(rw).Encode(fmt.Sprintf("[%s * ] match(%s) %v time\n", now.In(loc), testTime, time.Until(timestamps)))
		}

	}
}

func main() {

	http.HandleFunc("/makeMatch", makeMatch)

	http.HandleFunc("/unscheduledMatches", unschedueledMatches)

	http.HandleFunc("/scheduleMatches", scheduleMatches)

	fmt.Println("[*] listening on port 1337[*]")
	http.ListenAndServe(":1337", nil)
}
