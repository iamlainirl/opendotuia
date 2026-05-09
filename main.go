package main 

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

type Profile struct {
	Personaname string `json:"personaname"`
}

type Player struct {
	Profile Profile `json:"profile"`
}

type Match struct {
	MatchID int `json:"match_id"`
	RadiantWin bool `json:"radiant_win"`
	HeroID int `json:"hero_id"`
	Duration int `json:"duration"`
	GameMode int `json:"game_mode"`
	LobbyType int `json:"lobby_type"`
	Kills int `json:"kills"`
	Deaths int `json:"deaths"`
	Assists int `json:"assists"`
	//average_rank
	//leaver status
	//party type
	
}
func main() {
	var id string
	fmt.Scan(&id)
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s", id)
	resp, err := http.Get(url)
	matches_url := fmt.Sprintf("https://api.opendota.com/api/players/%s/recentMatches", id)
	resp2, err2 := http.Get(matches_url)
	if err != nil {
		panic(err)
	} 
	if err2 != nil {
		panic(err2)
	}
	body, _ := io.ReadAll(resp.Body)
	body2, _ := io.ReadAll(resp2.Body)
	var player Player
	var matches []Match
	json.Unmarshal(body, &player)
	fmt.Println(player.Profile.Personaname)
	json.Unmarshal(body2, &matches)
	if len(matches) == 0 {
		fmt.Println("This profile is private or no recent matches")
	}
	for _, match := range matches {
		fmt.Printf("%v | %d/%d/%d | %d:%02d\n", match.RadiantWin, match.Kills, match.Deaths, match.Assists, match.Duration/60, match.Duration%60)
	}
	return
}
