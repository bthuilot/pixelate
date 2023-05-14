FROM golang:1.20

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

# Start the Xorg server and run the application on container start
CMD /app/pixelate