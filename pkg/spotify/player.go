package spotify

import (
	"context"
)

func GetCurrentAlbumArtURL() (string, bool) {
	if client == nil {
		Init()
	}
	player, err := client.PlayerState(context.Background())

	if err != nil {
		return "", false
	}

	if !player.Playing {
		return "", false
	}

	images := player.Item.Album.Images

	if len(images) < 0 {
		return "", false
	}

	return images[0].URL, player.Playing
}