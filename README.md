# πxelate

[![Go Report Card](https://goreportcard.com/badge/github.com/bthuilot/pixelate)](https://goreportcard.com/report/github.com/bthuilot/pixelate)


*Note*: This project was inspired by [@rwardtech](https://www.tiktok.com/@rwardtech) on TikTok and his SpotiPi project which can be view [here on GitHub](https://github.com/ryanwa18/spotipi).


## Hardware Needed

1. [Raspberry Pi](https://www.raspberrypi.com/products/) & SD Card with Raspbian installed
2. [64x64 RGB Matrix](https://www.adafruit.com/product/3649)
3. [Adafruit Matrix Bonnet Hat](https://www.adafruit.com/product/3211)
4. [5V 10A A/C Adapter](https://www.adafruit.com/product/658)

## Up and Running

To Install:
1. navigate to the releases tab and download the latest GZip of the binary.
2. Create the file `config.yml` in the folder `/etc/pixelate/` (an example is shown in [`example.config.yml`](/example.config.yml))
3. Run the binary

### Build from source

1. Clone this repo onto the Raspberry Pi
2. Initialize Git Submodules (`git submodule update --init --recursive`)
3. run `make`
4. Create the `config.yml`
5. Run the binary


## Supported screens

A "Screen" represents a service that will render to the display.

The current list of "screen" values are:

- **Spotify**: Render the currently playing album to the board
- **Wifi QR Code**: Render a QR code to join a wifi network
