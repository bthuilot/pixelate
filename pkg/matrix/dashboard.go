package matrix

import (
	"SpotifyDash/pkg/server"
	"html/template"
	"net/http"
	"reflect"
)

const (
	SpotifyStateValue = "spotify"
	TickerStateValue = "ticker"
)

type configValue struct {
	FormattedName string
	Name string
	Value string
}

func Dashboard(w http.ResponseWriter, r* http.Request)  {
	t := template.Must(template.ParseFiles("./web/template/dashboard.tmpl"))
	w.WriteHeader(200)
	t.Execute(w, struct{
		CurrentState int
		StateTypes []string
	}{
		CurrentState: CurrentState,
		StateTypes: []string{SpotifyStateValue, TickerStateValue},
	})
}

func UpdateState(w http.ResponseWriter, r* http.Request)  {
	newState := r.PostForm.Get("state")
	switch newState {
	case SpotifyStateValue:
		CurrentState = SpotifyAlbum
	case TickerStateValue:
		CurrentState = Ticker
	default:
		server.TextResponse(w, "invalid state name", 404)
	}

	http.Redirect(w,r, "/", http.StatusSeeOther)
}

func parseConfig(value reflect.Value, t reflect.Type) []configValue {
	var types []configValue
	for i := 0; i < t.NumField(); i++{
		field := t.Field(i)
		types = append(types,
			configValue{
				FormattedName: field.Tag.Get("formattedName"),
				Name:          field.Tag.Get("formName"),
				Value:         reflect.Indirect(value).FieldByName(field.Name).String(),
			})
	}
	return types
}

func getCurrentConfigValues() map[string][]configValue {
	return map[string][]configValue{
		SpotifyStateValue: parseConfig(reflect.ValueOf(&spotifyConfig), reflect.TypeOf(spotifyConfig)),
		TickerStateValue: parseConfig(reflect.ValueOf(&tickerConfig), reflect.TypeOf(tickerConfig)),
	}
}

func CreateConfigPage(name string) {

}