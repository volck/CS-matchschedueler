package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)



type MatchConfig struct {
	Matchid              string `json:"matchid"`
	NumMaps              int    `json:"num_maps"`
	PlayersPerTeam       int    `json:"players_per_team"`
	MinPlayersToReady    int    `json:"min_players_to_ready"`
	MinSpectatorsToReady int    `json:"min_spectators_to_ready"`
	SkipVeto             bool   `json:"skip_veto"`
	VetoFirst            string `json:"veto_first"`
	SideType             string `json:"side_type"`
	Spectators           struct {
		Players []string `json:"players"`
	} `json:"spectators"`
	Maplist                []string `json:"maplist"`
	FavoredPercentageTeam1 int      `json:"favored_percentage_team1"`
	FavoredPercentageText  string   `json:"favored_percentage_text"`
	Team1                  struct {
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
	Matchtime time.Time `json:"matchtime"`
}

var matchConfigs []MatchConfig


func makeMatch(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

    var match MatchConfig
    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&match)
    if err != nil {
        json.NewEncoder(rw).Encode(err)
    } else {
		json.NewEncoder(rw).Encode(fmt.Sprintf("%s vs %s going to be played on %s schedueled for %s \n", match.Team1.Name, match.Team2.Name, match.Cvars.Hostname, match.Matchtime))
		fmt.Sprintf("%s vs %s going to be played on %s schedueled for %s \n", match.Team1.Name, match.Team2.Name, match.Cvars.Hostname, match.Matchtime)
		matchConfigs = append(matchConfigs, match)
	}

}


func getNumMatches(rw http.ResponseWriter, req *http.Request) {
	json.NewEncoder(rw).Encode(fmt.Sprintf("%d matches available", len(matchConfigs)))
}

func scheduleMatches(rw http.ResponseWriter, req *http.Request) {
	for match := range matchConfigs {
		fmt.Println(time.Until(matchConfigs[match].Matchtime))

	}
}



func main() {

http.HandleFunc("/makeMatch", makeMatch)

http.HandleFunc("/getNumMatches", getNumMatches)

http.HandleFunc("/scheduleMatches", scheduleMatches)


fmt.Println("[*] listening on port 1337[*]")
http.ListenAndServe(":1337", nil)
}
