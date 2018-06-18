FROM golang:1-alpine

ADD . /go/src/github/franxois/ggql-git/
WORKDIR /go/src/github/franxois/ggql-git/server/

RUN apk add --no-cache git

RUN go get -u github.com/vektah/gqlgen
RUN go generate graph/graph.go
RUN go get -u golang.org/x/vgo
RUN vgo build

RUN go get -u github.com/oxequa/realize