package config

import (
	"fmt"

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

type ServerConfig struct {
	Port        int
	Host        string
	ExternalURL string
}

type SpotifyConfig struct {
	ClientID     string `mapstructure:"clientID"`
	ClientSecret string `mapstructure:"clientSecret"`
}

type WifiQRConfig struct {
	SSID     string `mapstructure:"ssid"`
	AuthType string `mapstructure:"authType"`
	Password string `mapstructure:"password"`
}

type ConfigFile struct {
	Server     ServerConfig  `mapstructure:"server"`
	Spotify    SpotifyConfig `mapstructure:"spotify"`
	WifiQRCode WifiQRConfig  `mapstructure:"wifi_qrcode"`
	Logging    LogConfig     `mapstrucutre:"logging"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	LogFile   string `mapstructure:"log_file"`
	UseSTDOUT bool   `mapstructure:"use_stdout"`
}

// InitConfig will initialize the viper configuration
func InitConfig() (cfg ConfigFile, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/pixelate/")
	if err = viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("unable to read config file: %s", err)
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}
