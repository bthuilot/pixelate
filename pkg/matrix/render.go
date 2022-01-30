package matrix

import (
	"SpotifyDash/pkg/image_util"
	"SpotifyDash/pkg/spotify"
	"fmt"
	"log"
)


func RenderSpotify() {
	url, playing := spotify.GetCurrentAlbumArtURL()
	if !playing {
		RenderText("nothing currently playing")
	} else {
		RenderAlbum(url)
	}
}

func RenderText(text string) {
	fmt.Println(text)
	// TODO
}

func RenderAlbum(url string) {
	img, err := image_util.FromURL(url)
	if err != nil {
		log.Fatal(err)
	}
	thumbnail := image_util.Resize(img)
	fmt.Println(thumbnail.Bounds().String())
	// TODO
}