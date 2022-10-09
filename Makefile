include secrets.env
export

pixelate:
	cd src && go build -o ../pixelate main.go 

.PHONY: debug test install prod-vars release
debug:
	cd src && go run main.go

install: pixelate
	cp pixelate /usr/local/bin/pixelate

test:
	cd src && go test

prod-vars:
	export GIN_MODE=release

release: