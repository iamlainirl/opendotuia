package main 

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/charmbracelet/lipgloss"
)

type Profile struct {
	Personaname string `json:"personaname"`
}

type Player struct {
	Profile Profile `json:"profile"`
}

type Match struct {
	MatchID int `json:"match_id"`
	PlayerSlot int `json:"player_slot"`
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

var (
	winStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("76"))

	lossStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("1"))

	nameStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("7"))
)
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
	line := fmt.Sprintf(player.Profile.Personaname)
	fmt.Println(nameStyle.Render(line))
	json.Unmarshal(body2, &matches)
	if len(matches) == 0 {
		fmt.Println("This profile is private or no recent matches")
	}
	for _, match := range matches {
		win := (match.RadiantWin && match.PlayerSlot < 128) || (!match.RadiantWin && match.PlayerSlot >= 128)
		result := lossStyle.Render("LOSS")
		if win {
			result = winStyle.Render("WIN")
		}
		line2 := fmt.Sprintf("%v | %d/%d/%d | %d:%02d\n", result, match.Kills, match.Deaths, match.Assists, match.Duration/60, match.Duration%60)
		fmt.Println(line2)
	}
	return
}
