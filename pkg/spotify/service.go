package spotify

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/exp/rand"
	"image"
	"log"
	"net/http"
	"time"
)

const redirectURL = "http://matrix.thuilot.io:8080/spotify/callback"

var state = fmt.Sprintf("%d", rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Int63())

type Service struct {
	client *spotify.Client
	matrix chan image.Image
}

func (s *Service) Init(matrixChan chan image.Image, engine *gin.Engine) error {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
	)
	s.matrix = matrixChan

	clientChan := make(chan *spotify.Client)
	url := auth.AuthURL(state)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	engine.GET("/spotify/login", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, url)
	})
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
		tok, err := auth.Token(context, state, r)
		if err != nil {
			context.AbortWithStatus(http.StatusForbidden)
			log.Fatal(err)
			return
		}
		if st := r.FormValue("state"); st != state {
			context.AbortWithStatus(http.StatusNotFound)
			log.Fatal(fmt.Errorf("state mismatch: %s != %s\n", st, state))
			return
		}

		// use the token to get an authenticated client
		client := spotify.New(auth.Client(r.Context(), tok))
		_, _ = fmt.Fprintf(w, "Login Completed!")
		clientChannel <- client
	}
}
func (s *Service) Tick() error {
	if s.client == nil {
		s.matrix <- util.RenderText("Must login")
		return fmt.Errorf("please log into spotify")
	}
	img, err := s.RenderAlbumArt()
	if err != nil {
		s.matrix <- util.RenderText("error")
	} else {
		s.matrix <- img
	}
	return err
}

func (s *Service) GetConfig() api.ConfigStore {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetConfig(config api.ConfigStore) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (s *Service) RefreshDelay() time.Duration {
	return time.Second * 5
}

func (s *Service) GetID() string {
	return "spotify"
}
