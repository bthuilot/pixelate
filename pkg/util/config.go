package util

import (
	"github.com/spf13/viper"
)

const (
	// SpotifyClientID is the name of the Viper config key for the Spotify API client ID
	SpotifyClientID = "spotify_client_id"
	// SpotifyClientSecret is the name of the Viper config key for the Spotify API client secret
	SpotifyClientSecret = "spotify_client_secret"
	// ServerURL is the name of the Viper config key for the Server's URL
	ServerURL = "server_url"
)

// InitConfig will initialize the viper configuration
func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/pixelate/")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
