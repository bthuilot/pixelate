package matrix

import (
	"SpotifyDash/pkg/server"
	"html/template"
	"net/http"
	"strings"
)

const (
	SpotifyStateValue = "spotify"
	TickerStateValue  = "ticker"
)

type configValue struct {
	FormattedName string `json:"formatted_name"`
	Name          string `json:"form_name"`
	Value         string `json:"value"`
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./web/template/dashboard.tmpl"))
	w.WriteHeader(200)
	t.Execute(w, struct {
		CurrentState int
		StateTypes   []string
	}{
		CurrentState: CurrentState,
		StateTypes:   []string{SpotifyStateValue, TickerStateValue},
	})
}

func UpdateState(w http.ResponseWriter, r *http.Request) {
	newState := r.PostForm.Get("state")
	switch newState {
	case SpotifyStateValue:
		CurrentState = SpotifyAlbum
	case TickerStateValue:
		CurrentState = Ticker
	default:
		server.TextResponse(w, "invalid state name", 404)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ConfigSettings(w http.ResponseWriter, r *http.Request) {
	configPage := strings.TrimPrefix(r.URL.Path, "/config/")
	config, err := getConfig(configPage)
	if err != nil {
		server.TextResponse(w, `invalid request`, http.StatusBadRequest)
	}
	server.JsonResponse(w, config, http.StatusOK)
}
