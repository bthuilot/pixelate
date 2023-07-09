package screens

import (
	"context"
	"fmt"
	"github.com/bthuilot/pixelate/pkg/config"
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/bthuilot/pixelate/pkg/rendering"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"image"
	"math/rand"
	"net/http"
	"time"
)

type Spotify struct {
	authenticator *spotifyauth.Authenticator
	client        *spotify.Client
}

func (s Spotify) Render() (img image.Image, dur time.Duration, err error) {
	dur = time.Second * 15
	if s.client == nil {
		img = rendering.RenderText("go to homepage to login")
	}
	img, err = s.renderAlbumArt()
	return
}

func (s Spotify) GetConfig() map[string]string {
	return nil
}

func (s Spotify) UpdateConfig(m map[string]string) error {
	return nil
}

func (s Spotify) GetHTMLPage() (attrs []display.HTMLAttributes) {
	attrs = append(attrs, display.HTMLLink{
		Text: "Login with Spotify",
		Href: s.authenticator.AuthURL(state),
	})
	if s.client != nil {
		var content string
		if usr, err := s.client.CurrentUser(context.Background()); err != nil {
			content = "unable to retrieve user"
		} else {
			content = fmt.Sprintf("%s <%s>", usr.DisplayName, usr.Email)
		}
		attrs = append(attrs, display.HTMLText{
			Content: fmt.Sprintf("Currently logged in as: %s", content),
		})
	}
	return
}

func (s Spotify) Init(r *gin.RouterGroup) (name string, err error) {
	name = "Spotify"
	r.GET("/screens/spotify/callback", s.authCallback)
	return
}

func (s Spotify) authCallback(c *gin.Context) {
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

var state = fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63())

func (s Spotify) renderAlbumArt() (img image.Image, err error) {
	if s.client == nil {
		return rendering.RenderText("please sign in on homepage"), nil
	}
	player, err := s.client.PlayerState(context.Background())
	if err != nil {
		return nil, err
	}
	if !player.Playing {
		return rendering.RenderText("No songs playing"), nil
	}

	images := player.Item.Album.Images

	if len(images) > 0 {
		url := images[0].URL
		img, err = rendering.ImageFromURL(url)
		return imaging.Resize(img, 64, 64, imaging.Lanczos), err
	}
	return nil, fmt.Errorf("no album art images returned from API")
}

func NewSpotifyScreen(cfg config.ConfigFile) display.Screen {
	return Spotify{
		authenticator: spotifyauth.New(
			spotifyauth.WithRedirectURL(fmt.Sprintf("%s/api/spotify/callback", cfg.Server.ExternalURL)),
			spotifyauth.WithScopes(
				spotifyauth.ScopeUserReadCurrentlyPlaying,
				spotifyauth.ScopeUserReadPlaybackState,
				spotifyauth.ScopeUserReadEmail),
			spotifyauth.WithClientID(cfg.Spotify.ClientID),
			spotifyauth.WithClientSecret(cfg.Spotify.ClientSecret),
		),
		client: nil,
	}
}
