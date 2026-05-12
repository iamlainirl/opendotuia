package main 

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"encoding/json"
	"github.com/charmbracelet/lipgloss"
	_ "embed"
)

type Profile struct {
	Personaname string `json:"personaname"`
}

type Player struct {
	Profile Profile `json:"profile"`
}

//go:embed heroes.json
var heroesByte []byte


type Hero struct{
	ID int `json:"id"`
	LocalizedName string `json:"localized_name"`
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

	heroStyle = lipgloss.NewStyle().
	Width(20).
	Foreground(lipgloss.Color("50"))

	kdaStyle = lipgloss.NewStyle().
	Width(10).
	Align(lipgloss.Center)
)

func main() {
	var heroesList []Hero
	json.Unmarshal(heroesByte, &heroesList)

	heroMap := make(map[int]string)
	for _, h := range heroesList {
		heroMap[h.ID] = h.LocalizedName
	}

	var id string
	fmt.Scan(&id)
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s", id)
	resp, err := http.Get(url)
	matches_url := fmt.Sprintf("https://api.opendota.com/api/players/%s/recentMatches", id)
	resp2, err2 := http.Get(matches_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	body2, _ := io.ReadAll(resp2.Body)
	var player Player
	var matches []Match

	json.Unmarshal(body, &player)
	json.Unmarshal(body2, &matches)
	if len(matches) == 0 {
		fmt.Println("This profile is private or no recent matches")
	}

	fmt.Println(nameStyle.Render("\nPlayer : " + player.Profile.Personaname))
	fmt.Println(strings.Repeat("--", 25))

	for _, match := range matches {
		win := (match.RadiantWin && match.PlayerSlot < 128) || (!match.RadiantWin && match.PlayerSlot >= 128)
		hName, _ := heroMap[match.HeroID]
		kda := fmt.Sprintf("%d/%d/%d", match.Kills, match.Deaths, match.Assists)

		result := lossStyle.Render("LOSS")
		if win {
			result = winStyle.Render("WIN")
		}
		fmt.Printf("%s | %s | %s | %d:%02d\n",
		result,
		heroStyle.Render(hName),
		kdaStyle.Render(kda),
		match.Duration/60, match.Duration%60,
	)
}

	return
}
