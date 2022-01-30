package matrix

import "net/url"

var tickerConfig = struct {
	tickerSymbol string `default:"BX" formName:"ticker_symbol" formattedName:"Ticker Symbol"`
}{}

func updateTickerConfig(newValues url.Values) {
	tickerConfig.tickerSymbol = newValues.Get("ticker_symbol")
}

var spotifyConfig = struct {
}{}

func updateSpotifyConfig(newValues url.Values) {
}