package services

import (
	"SpotifyDash/pkg/util"
	"context"
	"fmt"
	"image"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Spotify struct {
	client *spotify.Client
	matrix chan image.Image
}

func (s Spotify) GetName() ID {
	return "Spotify"
}

func (s Spotify) Run(_ Config) chan Command {
	// TODO
	cmdChannel := make(chan Command)
	go func() {
	exit_routine:
		for {
			select {
			case cmd := <-cmdChannel:
				code := cmd.Code
				switch code {
				case Stop:
					break exit_routine
				case Tick:
					s.tick()
				case Update:
					// No config so don't matter
				}
			}
		}

	}()
	return cmdChannel
}

func (s Spotify) GetDefaultConfig() Config {
	return Config{}
}

var state = fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63())

func (s Spotify) Init(matrixChan chan image.Image) (page SetupPage) {
	baseURL := "localhost:8080"
	if newBaseUrl := os.Getenv("SPOTIFY_CALLBACK_URL"); newBaseUrl != "" {
		baseURL = newBaseUrl
	}
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(fmt.Sprintf("http://%s/spotify/callback", baseURL)),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
	)
	s.matrix = matrixChan
	page = append(page, Button{
		Link: auth.AuthURL(state),
		Name: "Login with Spotify",
	})
	go func() {
		m := http.NewServeMux()
		svr := http.Server{Addr: ":7000", Handler: m}
		m.HandleFunc("/spotify/callback", s.spotifyAuthCallback(auth, svr))
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		log.Printf("Finished")
	}()
	return
}

func (s Spotify) spotifyAuthCallback(auth *spotifyauth.Authenticator,
	svr http.Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(r.Context(), state, r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("error"))
			log.Fatal(err)
			return
		}
		if st := r.FormValue("state"); st != state {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("error"))
			log.Fatal(fmt.Errorf("state mismatch: %s != %s\n", st, state))
			return
		}

		// use the token to get an authenticated client
		s.client = spotify.New(auth.Client(r.Context(), tok))
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Add("Location", "matrix.thuilot.io") // Go back to main server
		svr.Shutdown(context.Background())
	}
}

func (s Spotify) tick() error {
	if s.client == nil {
		s.matrix <- util.RenderText("go to /setup to login")
	}
	img, err := s.renderAlbumArt()
	if err != nil {
		s.matrix <- util.RenderText("error rendering album art")
	} else {
		s.matrix <- img
	}
	return err
}

func (s Spotify) GetTickInterval() time.Duration {
	return time.Minute
}

func (s *Spotify) renderAlbumArt() (img image.Image, err error) {
	player, err := s.client.PlayerState(context.Background())
	if err != nil {
		return nil, err
	}
	if !player.Playing {
		return util.RenderText("No songs playing"), nil
	}

	images := player.Item.Album.Images

	if len(images) > 0 {
		url := images[0].URL
		img, err := util.FromURL(url)
		if err != nil {
			return nil, err
		}
		thumbnail := util.Resize(img)
		return thumbnail, nil
	}
	return nil, fmt.Errorf("no album art images returned from API")
}
