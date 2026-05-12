package main 

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"encoding/json"
	"github.com/charmbracelet/lipgloss"
	_ "embed"

	tea "charm.land/bubbletea/v2"
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
	AverageRank int `json:"average_rank"`
	LeaverStatus int `json:"leaver_status"`
	PartySize *int `json:"party_size"`
}

// screen 



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

	abandonedStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("166"))
)

type model struct {
	input string
	screen string
	cursor int
	player Player
	matches []Match
	heroes map[int]Hero
	err error
	//width int 
	//height int

}

func loadHeroes() map[int]Hero {
	var heroes []Hero
	json.Unmarshal(heroesByte, &heroes)
	res := make(map[int]Hero)
	for _, hero := range heroes {
		res[hero.ID] = hero
	}
	return res
}

func initialModel() model {
	return model {
		screen: "input",
		heroes: loadHeroes(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type playerMsg Player
type matchesMsg []Match

func fetchPlayer(id string) tea.Cmd {
	return func() tea.Msg {
// возможно объявитьчтение строки?
		url := fmt.Sprintf("https://api.opendota.com/api/players/%s", id)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		
		var player Player
		json.Unmarshal(body, &player)
		return playerMsg(player)
	}
}

func fetchMatches(id string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.opendota.com/api/players/%s/recentMatches", id)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
			//добавить обработку ошибок?
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var matches []Match
		json.Unmarshal(body, &matches)
		return matchesMsg(matches)
	}
}

func Rank(r int) string {
	star := r % 10
	switch {
	case r < 20: return fmt.Sprintf("Herald %d", star)
	case r < 30: return fmt.Sprintf("Guardian %d", star)
	case r < 40: return fmt.Sprintf("Crusader %d", star)
	case r < 50: return fmt.Sprintf("Archon %d", star)
	case r < 60: return fmt.Sprintf("Legend %d", star)
	case r < 70: return fmt.Sprintf("Ancient %d", star)
	case r < 80: return fmt.Sprintf("Divine %d", star)
	case r >= 80: return "Immortal"
	default: return "Unknown"
	}
}

func partyString(p *int) string {
	if p == nil {
		return "Solo"
	}
	if *p == 1 {
		return "Solo"
	}
	return fmt.Sprintf("Party %d", *p)
}

func leaverStr(l int) string {
	if l == 0 {
		return ""
	}
	return "Abandoned"
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String(){

		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "enter":
			m.screen = "loading"
			return m, tea.Batch(fetchPlayer(m.input), fetchMatches(m.input))

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			m.input = m.input + msg.String()
		}

	case playerMsg:
		m.player = Player(msg)
		return m, nil
	
	case matchesMsg:
		m.matches = []Match(msg)
		m.screen = "matches"
		return m, nil
	}
	return m, nil
}
func (m model) View() tea.View {
	if m.screen == "input" {
		return tea.NewView("Enter player ID:\n> " + m.input)
	}
	if m.screen == "loading" {
		return tea.NewView("Loading...")
	} else {
		s := nameStyle.Render(m.player.Profile.Personaname) + "\n"
		for _, match := range m.matches {
			win := (match.RadiantWin && match.PlayerSlot < 128) || (!match.RadiantWin && match.PlayerSlot >= 128) 
			heroName := m.heroes[match.HeroID].LocalizedName
			kda := fmt.Sprintf("%d/%d/%d", match.Kills, match.Deaths, match.Assists)
			result := lossStyle.Render("LOSS")
			if win {
				result = winStyle.Render("WIN")
			}
			if match.LeaverStatus > 0 {
				result = abandonedStyle.Render("Abandoned")
			}
			line := fmt.Sprintf("%s | %s | %s | %s | %s | %d:%02d",
				result,
				heroStyle.Render(heroName),
				partyString(match.PartySize),
				Rank(match.AverageRank),
				kdaStyle.Render(kda),
				match.Duration/60, match.Duration%60,
			)
			s += line + "\n"
		}
		return tea.NewView(s)
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return
}
