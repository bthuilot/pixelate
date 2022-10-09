include secrets.env
export

pixelate:
	cd src && go build main.go -o ../

.PHONY: debug test install
debug:
	cd src && go run main.go

install: pixelate
	cp pixelate /usr/local/bin/pixelate

test:
	cd src && go test


