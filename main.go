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

func main() {
	var id string
	fmt.Scan(&id)
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		var player Player
		json.Unmarshal(body, &player)
		fmt.Println(player.Profile.Personaname)
	}
	return
}
