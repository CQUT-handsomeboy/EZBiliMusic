FROM golang:1.24-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN go build -v -o ./app .

CMD ["./app"]
