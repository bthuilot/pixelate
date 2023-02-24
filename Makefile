SOURCEDIR = pkg
# Go Sources
API_GO_SOURCES = $(wildcard $(SOURCEDIR)/api/*.go)
HTTP_GO_SOURCES = $(wildcard $(SOURCEDIR)/http/*.go)
MATRIX_GO_SOURCES = $(wildcard $(SOURCEDIR)/matrix/*.go)
RENDERING_GO_SOURCES = $(wildcard $(SOURCEDIR)/rendering/*.go)
UTIL_GO_SOURCES = $(wildcard $(SOURCEDIR)/util/*.go)
GO_SOURCES = $(API_GO_SOURCES) $(HTTP_GO_SOURCES) $(MATRIX_GO_SOURCES) $(RENDERING_GO_SOURCES) $(UTIL_GO_SOURCES)

ASSEST_DIR = assets
# Embedded sources
TEMPLATE_SOURCES = $(wildcard $(ASSEST_DIR)/web/templates/*.tmpl)
STATIC_SOURCES = $(wildcard $(ASSEST_DIR)/web/static/*)
FONTS = $(wildcard $(ASSEST_DIR)/fonts/*.ttf)

pixelate: rpi-lib $(GO_SOURCES) $(TEMPLATE_SOURCES) $(STATIC_SOURCES) $(FONTS)
	go build -o pixelate

.PHONY: clean test install prod-vars release rpi-lib
rpi-lib:
	$(MAKE) -C third_party/rpi-rgb-led-matrix/lib all

clean:
	rm -f pixelate

install: pixelate
	cp pixelate /usr/local/bin/pixelate

test:
	go test

prod-vars:
	export GIN_MODE=release

release: clean prod-vars pixelate