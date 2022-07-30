package spotify

// func renderNotPlaying() image.Image {
// 	return
// }

// func (s *Service) RenderAlbumArt() (img image.Image, err error) {
// 	player, err := s.client.PlayerState(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !player.Playing {
// 		return renderNotPlaying(), nil
// 	}

// 	images := player.Item.Album.Images

// 	if len(images) > 0 {
// 		url := images[0].URL
// 		img, err := util.FromURL(url)
// 		if err != nil {
// 			return nil, err
// 		}
// 		thumbnail := util.Resize(img)
// 		return thumbnail, nil
// 	}
// 	return util.RenderText("Loading..."), nil
// }
