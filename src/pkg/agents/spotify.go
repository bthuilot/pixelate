package agents

import (
	"context"
	"fmt"
	"github.com/bthuilot/pixelate/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	"image"
	"math/rand"
	"net/http"
	"os"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Spotify struct {
	authenticator *spotifyauth.Authenticator
	client        *spotify.Client
	cfg           Config
}

func (s *Spotify) SetConfig(config Config) error {
	s.cfg = config
	return nil
}

func (s *Spotify) RegisterEndpoints(r *gin.Engine) {
	s.authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(fmt.Sprintf("http://%s/spotify/callback", os.Getenv("MATRIX_SERVER_URL"))),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
	)
	r.Any("/spotify/callback", s.authCallback)
}

func (s *Spotify) authCallback(c *gin.Context) {
	tok, err := s.authenticator.Token(c.Request.Context(), state, c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	if recvState, exist := c.GetQuery("state"); !exist || recvState != state {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "state mismatch",
		})
		return
	}

	// use the token to get an authenticated client
	s.client = spotify.New(s.authenticator.Client(c.Request.Context(), tok))
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Spotify) GetTick() time.Duration {
	return time.Second * 15
}

func (s *Spotify) GetConfig() Config {
	return s.cfg
}

func (s *Spotify) GetAdditionalHTML() []Attribute {
	return []Attribute{
		Button{
			Name: "Login with Spotify",
			Link: s.authenticator.AuthURL(state),
		},
	}
}

func (s *Spotify) Render(img chan image.Image) {
	if s.client == nil {
		img <- util.RenderText("go to homepage to login")
	}
	if albumArt, err := s.renderAlbumArt(); err != nil {
		img <- util.RenderText("error rendering album art")
	} else {
		img <- albumArt
	}
}

func NewSpotify() Renderer {
	return &Spotify{
		authenticator: nil,
		client:        nil,
		cfg:           nil,
	}
}

func (s *Spotify) GetName() ID {
	return "Spotify"
}

var state = fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63())

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
