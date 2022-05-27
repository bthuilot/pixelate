package spotify

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/image_util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"image"
	"log"
	"net/http"
)

type Service struct {
	client *spotify.Client
	matrix chan image.Image
}

func (s Service) Init(matrixChan chan image.Image, engine *gin.Engine) error {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(RedirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
	)
	s.matrix = matrixChan

	clientChan := make(chan *spotify.Client)
	url := auth.AuthURL(State)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	engine.GET("/spotify/callback", createCallback(clientChan, auth))
	go func() {
		select {
		case s.client = <-clientChan:
			return
		}
	}()
	return nil
}

func createCallback(clientChannel chan *spotify.Client, auth *spotifyauth.Authenticator) gin.HandlerFunc {
	return func(context *gin.Context) {
		r, w := context.Request, context.Writer
		tok, err := auth.Token(context, State, r)
		if err != nil {
			context.AbortWithStatus(http.StatusForbidden)
			log.Fatal(err)
			return
		}
		if st := r.FormValue("state"); st != State {
			context.AbortWithStatus(http.StatusNotFound)
			log.Fatal(fmt.Errorf("state mismatch: %s != %s\n", st, State))
			return
		}

		// use the token to get an authenticated client
		client := spotify.New(auth.Client(r.Context(), tok))
		_, _ = fmt.Fprintf(w, "Login Completed!")
		clientChannel <- client
	}
}
func (s Service) Tick() {
	if s.client == nil {
		return
	}

	player, _ := s.client.PlayerState(context.Background())

	if !player.Playing {
		return
	}

	images := player.Item.Album.Images

	if len(images) > 0 {
		url := images[0].URL
		img, err := image_util.FromURL(url)
		if err != nil {
			log.Fatal(err)
		}
		thumbnail := image_util.Resize(img)
		s.matrix <- thumbnail
	}
}

func (s Service) GetConfig() api.ConfigStore {
	//TODO implement me
	panic("implement me")
}

func (s Service) SetConfig(config api.ConfigStore) error {
	//TODO implement me
	panic("implement me")
	return nil
}
