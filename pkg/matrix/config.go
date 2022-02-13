package matrix

import (
	"fmt"
	"net/url"
	"reflect"
)

var tickerConfig = struct {
	tickerSymbol string `formName:"ticker_symbol" formattedName:"Ticker Symbol"`
}{
	tickerSymbol: "BX",
}

func updateTickerConfig(newValues url.Values) {
	tickerConfig.tickerSymbol = newValues.Get("ticker_symbol")
}

var spotifyConfig = struct {
}{}

func updateSpotifyConfig(newValues url.Values) {
}

func getConfig(name string) ([]configValue, error) {
	switch name {
	case SpotifyStateValue:
		return parseConfig(&spotifyConfig, spotifyConfig), nil
	case TickerStateValue:
		return parseConfig(&tickerConfig, tickerConfig), nil
	default:
		return nil, fmt.Errorf("no config with name %s", name)
	}
}

func parseConfig(ptr interface{}, config interface{}) []configValue {
	value, t := reflect.ValueOf(ptr), reflect.TypeOf(config)
	var types []configValue
	for i := 0; i < t.NumField(); i++ {
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
