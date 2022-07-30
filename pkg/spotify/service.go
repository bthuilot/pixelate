package spotify

// var state = fmt.Sprintf("%d", rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Int63())

// type Service struct {
// 	matrix chan image.Image
// 	url    string
// }

// func (s *Service) Init(matrixChan chan image.Image, engine *gin.Engine) error {
// 	baseURL := "localhost:8080"
// 	if newBaseUrl := os.Getenv("SPOTIFY_CALLBACK_URL"); newBaseUrl != "" {
// 		baseURL = newBaseUrl
// 	}
// 	auth := spotifyauth.New(
// 		spotifyauth.WithRedirectURL(fmt.Sprintf("http://%s/spotify/callback", baseURL)),
// 		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
// 	)
// 	s.matrix = matrixChan

// 	clientChan := make(chan *spotify.Client)
// 	s.url = auth.AuthURL(state)

// 	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", s.url)
// 	engine.GET("/spotify/login", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "spotify_login.tmpl", s.url)
// 	})
// 	engine.GET("/spotify/generate", func(c *gin.Context) {
// 		baseURL := "localhost:8080"
// 		if newBaseUrl := os.Getenv("SPOTIFY_CALLBACK_URL"); newBaseUrl != "" {
// 			baseURL = newBaseUrl
// 		}
// 		auth := spotifyauth.New(
// 			spotifyauth.WithRedirectURL(fmt.Sprintf("http://%s/spotify/callback", baseURL)),
// 			spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
// 		)
// 		s.url = auth.AuthURL(state)
// 	})
// 	engine.GET("/spotify/callback", createCallback(clientChan, auth))
// 	go func() {
// 		select {
// 		case s.client = <-clientChan:
// 			return
// 		}
// 	}()
// 	return nil
// }

// func (s *Service) Tick() error {

// }

// func (s *Service) GetConfig() api.ConfigStore {
// 	return api.ConfigStore{}
// }

// func (s *Service) SetConfig(config api.ConfigStore) error {
// 	return nil
// }

// func (s *Service) RefreshDelay() time.Duration {
// 	return time.Second * 5
// }

// func (s *Service) GetID() string {
// 	return "spotify"
// }
