FROM golang:1.21-bookworm

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
ENV DISPLAY :0
EXPOSE 6000

CMD /app/pixelate