package spotify

import (
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
)

const RedirectURL = "http://localhost:8080/spotify/callback"
const State = "test!"

var client *spotify.Client = nil

var ch = make(chan *spotify.Client)

func Init() {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(RedirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState),
		)

	url := auth.AuthURL(State)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	http.HandleFunc("/spotify/callback", createCallback(auth))
	client = <-ch
}

func createCallback(auth *spotifyauth.Authenticator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r * http.Request) {
		tok, err := auth.Token(r.Context(), State, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != State {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, State)
		}

		// use the token to get an authenticated client
		client := spotify.New(auth.Client(r.Context(), tok))
		fmt.Fprintf(w, "Login Completed!")
		ch <- client
	}
}