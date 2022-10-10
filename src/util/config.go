package util

import (
	"github.com/spf13/viper"
)

const (
	SpotifyClientID     = "spotify_client_id"
	SpotifyClientSecret = "spotify_client_secret"
	ServerURL           = "server_url"
)

func InitConfig() error {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("yaml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/pixelate/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.pixelate") // call multiple times to add many search paths
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	return viper.ReadInConfig()            // Find and read the config file
}
