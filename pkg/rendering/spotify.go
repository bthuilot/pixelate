package rendering

import (
	"context"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"time"

	"github.com/disintegration/imaging"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Spotify struct {
	clientID      string
	clientSecret  string
	serverURL     string
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
		spotifyauth.WithRedirectURL(fmt.Sprintf("http://%s/spotify/callback", s.serverURL)),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserReadEmail),
		spotifyauth.WithClientID(s.clientID),
		spotifyauth.WithClientSecret(s.clientSecret),
	)
	r.GET("/spotify/callback", s.authCallback)
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

func (s *Spotify) GetAdditionalConfig() []ConfigAttribute {
	attrs := []ConfigAttribute{
		Link{
			Name: "Login with Spotify",
			Href: s.authenticator.AuthURL(state),
		},
	}
	if s.client != nil {
		var content string
		if usr, err := s.client.CurrentUser(context.Background()); err != nil {
			content = "unable to retrieve user"
		} else {
			content = fmt.Sprintf("%s <%s>", usr.DisplayName, usr.Email)
		}
		attrs = append(attrs, Text{
			Content: fmt.Sprintf("Currently logged in as: %s", content),
		})
	}
	return attrs
}

func (s *Spotify) NextFrame() (img image.Image) {
	if s.client == nil {
		return RenderText("go to homepage to login")
	}
	if albumArt, err := s.renderAlbumArt(); err != nil {
		return RenderText("error rendering album art")
	} else {
		return albumArt
	}
}

func NewSpotifyAgent(clientID, clientSecret, serverURL string) Agent {
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
	if s.client == nil {
		return RenderText("please sign in on homepage"), nil
	}
	player, err := s.client.PlayerState(context.Background())
	if err != nil {
		return nil, err
	}
	if !player.Playing {
		return RenderText("No songs playing"), nil
	}

	images := player.Item.Album.Images

	if len(images) > 0 {
		url := images[0].URL
		img, err = ImageFromURL(url)
		return imaging.Resize(img, 64, 64, imaging.Lanczos), err
	}
	return nil, fmt.Errorf("no album art images returned from API")
}
