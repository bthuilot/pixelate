package screens

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/bthuilot/pixelate/pkg/config"
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type WifiQRCode struct {
	authType string
	ssid     string
	password string
}

func (w *WifiQRCode) Render() (img image.Image, dur time.Duration, err error) {
	var png []byte
	dur = -1
	wifiQRStr := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", w.authType, w.ssid, w.password)
	if png, err = qrcode.Encode(wifiQRStr, qrcode.Medium, 64); err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(png))
	return
}

func (w *WifiQRCode) GetConfig() map[string]string {
	return nil
}

func (w *WifiQRCode) UpdateConfig(m map[string]string) error {
	return nil
}

func (w *WifiQRCode) GetHTMLPage() (attrs []display.HTMLAttribute) {
	return
}

func (w *WifiQRCode) Init(r *gin.RouterGroup) (name string, err error) {
	name = "Wifi QR Code"
	return
}

func NewWifiQRCode(cfg config.ConfigFile) display.Screen {
	return &WifiQRCode{
		authType: cfg.WifiQRCode.AuthType,
		ssid:     cfg.WifiQRCode.SSID,
		password: cfg.WifiQRCode.Password,
	}
}
