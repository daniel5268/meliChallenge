FROM golang:1.18 as builder

WORKDIR /ws

RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get install -y migrate

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
