package matrix

import (
	"fmt"
	"net/http"
	"time"
)

const (
	SpotifyAlbum = iota
	Ticker
)

var CurrentState = SpotifyAlbum

func Init() {
	http.HandleFunc("/", Dashboard)
	http.HandleFunc("/state", UpdateState)
	for {
		switch CurrentState {
		case SpotifyAlbum:
			RenderSpotify()
		case Ticker:
			RenderText(fmt.Sprintf("%s", tickerConfig.tickerSymbol))
		}
		time.Sleep(getWaitTime(CurrentState))
	}
	// TODO
}


func getWaitTime(state int) time.Duration {
	switch state {
	case SpotifyAlbum:
		return 60
	case Ticker:
		return 120
	default:
		return 360
	}
}