package ticker

import (
	"SpotifyDash/pkg/util"
	"os"
)

const apiURL = "https://www.alphavantage.co/query"

var apiKey = os.Getenv("ALPHA_VANTAGE_API_KEY")

type StockInfo struct {
	Ticker    string
	Price     float64
	dayChange float64
}

type TimeSeriesPoint struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type GlobalQuote struct {
	Symbol        string `json:"01. symbol"`
	Open          string `json:"02. open"`
	High          string `json:"03. high"`
	Low           string `json:"04. low"`
	Price         string `json:"05. price"`
	Volume        string `json:"06. volume"`
	PreviousClose string `json:"08. previous close"`
	Change        string `json:"09. change"`
	ChangePercent string `json:"10. change percent"`
}

type GlobalQuoteResponse struct {
	Quote GlobalQuote `json:"Global Quote"`
}

type StockResponse struct {
	MetaData   map[string]string          `json:"MetaData"`
	TimeSeries map[string]TimeSeriesPoint `json:"Time Series (60min)"`
}

func getStockInfo(ticker string) (quote GlobalQuote, err error) {
	var response GlobalQuoteResponse
	params := map[string]string{
		"symbol":   ticker,
		"apikey":   apiKey,
		"function": "GLOBAL_QUOTE",
	}
	err = util.HTTPRequest[GlobalQuoteResponse](apiURL, params, nil, nil, &response)
	quote = response.Quote
	return
}