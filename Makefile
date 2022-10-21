SOURCEDIR = src
# Go Sources
GO_SOURCES = $(wildcard $(SOURCEDIR)/*.go)
# Embded sources
HTML_SOURCES = $(wildcard $(SOURCEDIR)/web/templates/*.tmpl)
JS_SOURCES = $(wildcard $(SOURCEDIR)/web/static/js/*.js)
CSS_SOURCES = $(wildcard $(SOURCEDIR)/web/static/css/*.css)

pixelate: $(GO_SOURCES) $(HTML_SOURCES) $(JS_SOURCES) $(CSS_SOURCES) 
	go build -o ../pixelate

.PHONY: clean debug test install prod-vars release
clean:
	rm -f pixelate

debug: export MATRIX_EMULATOR = 1

debug:
	go run main.go

install: src/pixelate
	cp pixelate /usr/local/bin/pixelate

test:
	cd src && go test

prod-vars:
	export GIN_MODE=release

release: clean prod-vars pixelate