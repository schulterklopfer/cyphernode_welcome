FROM golang:1.12.1-alpine3.9

RUN apk add git bash curl ca-certificates

RUN mkdir -p $GOPATH/src/cyphernode_status
RUN mkdir -p /data

ADD cnAuth $GOPATH/src/cyphernode_status/cnAuth
ADD static $GOPATH/src/cyphernode_status/static
ADD templates $GOPATH/src/cyphernode_status/templates
COPY main.go $GOPATH/src/cyphernode_status

WORKDIR $GOPATH/src/cyphernode_status

RUN go get

RUN go build main.go
RUN chmod +x $GOPATH/src/cyphernode_status/main

ENV PATH=$PATH:$GOPATH/src/cyphernode_status/

CMD ["main"]