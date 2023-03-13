FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build cmd/daemon/main.go

WORKDIR /dist

RUN cp /build/main .

RUN cp /build/config.cfg .

EXPOSE 8881/tcp

CMD ["/dist/main"]
