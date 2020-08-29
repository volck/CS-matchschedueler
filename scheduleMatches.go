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
		if !matchIsStale(match) {
			json.NewEncoder(rw).Encode(fmt.Sprintf("%s vs %s going to be played on %s schedueled for %s", match.Team1.Name, match.Team2.Name, match.Cvars.Hostname, match.Matchtime))
			fmt.Sprintf("%s vs %s going to be played on %s schedueled for %s", match.Team1.Name, match.Team2.Name, match.Cvars.Hostname, match.Matchtime)

			matchConfigs = append(matchConfigs, match)
		} else {
			json.NewEncoder(rw).Encode(fmt.Sprintf("%s vs %s which was going to be played on %s is stale, and will not be added to matchconfigs ", match.Team1.Name, match.Team2.Name,match.Matchtime))

		}


	}
}

func matchIsStale(match MatchConfig)(stale bool) {
		when, err := getMatchTimes(match)
		if err != nil {
			fmt.Println(err)
		}
		if when < time.Hour {
			if when < 0 {
				return true
			} else {
				return false
			}
	}
	return false
}


func getMatchTimes(match MatchConfig)(when time.Duration, err error) {
		layout := "02-01-2006 15:04:05"
		testTime := match.Matchtime
		parsedTime, err := time.Parse(layout, testTime)
		loc, _ := time.LoadLocation("Europe/Oslo")
		timestamps := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), loc)
		if err != nil {
			return
		}

		when = time.Until(timestamps)
		return
}



func scheduleMatches(rw http.ResponseWriter, req *http.Request) {
	if len(matchConfigs) <= 0 {
		json.NewEncoder(rw).Encode("[ * ] no matches to schedule")
	}
	for match := range matchConfigs {
		matchtime, err := getMatchTimes(matchConfigs[match])
		if err != nil {
			fmt.Println(err)
		}
		if matchtime < time.Hour {
			json.NewEncoder(rw).Encode(fmt.Sprintf("%s vs %s needs to be scheduled now(%v)", matchConfigs[match].Team1,matchConfigs[match].Team2, matchtime))
 		} else {
 			json.NewEncoder(rw).Encode(fmt.Sprintf("%s vs %s (%v) is probably a long way away so we wont worry. yet.", matchConfigs[match].Team1.Name,matchConfigs[match].Team2.Name, matchtime))
		}
	}


}

func main() {

	http.HandleFunc("/makeMatch", makeMatch)


	http.HandleFunc("/scheduleMatches", scheduleMatches)

	fmt.Println("[*] listening on port 1337[*]")
	http.ListenAndServe(":1337", nil)
}
