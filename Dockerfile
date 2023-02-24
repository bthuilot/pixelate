# This Dockerfile serves as a X11 forwarder for running the
# application on systems where the rpi-rgb-led-matrix library cannot be built
FROM ubuntu:latest

RUN apt-get update && apt-get install -y xauth xorg openbox curl make build-essential

ENV PATH="$PATH:/usr/local/go/bin"
ENV GO_VERSION="1.20.1"
# Install Golang
RUN curl -sL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -xvz -C /usr/local

WORKDIR /app
COPY main.go /app/
COPY go.mod /app/
COPY go.sum /app/
COPY Makefile /app/
COPY config.yml /app/
COPY pkg/ /app/pkg/
COPY assets/ /app/assets/
COPY third_party/ /app/third_party/

RUN make release
ENV MATRIX_EMULATOR=1
EXPOSE 6000
ENTRYPOINT ["/app/pixelate"]