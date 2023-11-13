FROM golang:1.21 as builder

# ENV CGO_ENABLED 0

#Install Dependencies
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y ffmpeg
RUN go install github.com/bwmarrin/dca/cmd/dca@latest

WORKDIR /src

ADD . /src

RUN go build -o /app ./cmd/main.go

EXPOSE 8080

CMD ["/app"]