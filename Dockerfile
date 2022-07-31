FROM golang:1.18-bullseye

# RUN apt update && apt install gcc

WORKDIR /app

ENV GO111MODULE="on"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY pkg ./pkg/
COPY lib ./lib/
COPY main.go ./

WORKDIR /app/lib/rpi-rgb-led-matrix

# TODO make env

# TODO fix this, either add a make in the root,
#  or have go do it (first by figuring out how thats done).
RUN make

WORKDIR /app

RUN go build -o /pixelate -v

EXPOSE 8080 # WEb server
EXPOSE 7000 # Spotify call back

CMD ["/pixelate"]

