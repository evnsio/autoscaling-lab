FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/evns/sqrt-server/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/server
ENTRYPOINT ["/go/bin/server"]