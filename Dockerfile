FROM golang:1.24-alpine

WORKDIR /usr/src/app

RUN go install github.com/charmbracelet/skate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o ./app .

CMD ["./app"]
