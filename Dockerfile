# This Dockerfile serves as a X11 forwarder for running the
# application on systems where the rpi-rgb-led-matrix library cannot be built
#FROM ubuntu:latest
#
#ENV TZ=America/New_York
#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
FROM ubuntu:latest

# Install necessary packages for X11 forwarding and OpenGL
RUN apt-get update && apt-get install -y \
    curl make build-essential \
    xauth \
    xserver-xorg-core \
    x11-apps \
    mesa-utils \
    libgl1-mesa-glx \
    libgl1-mesa-dri \
    x11-xserver-utils


ENV PATH="$PATH:/usr/local/go/bin"
ENV GO_VERSION="1.20.1"
# Install Golang
RUN curl -sL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -xz -C /usr/local

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