FROM golang:1.10-alpine3.8
RUN apk add --no-cache git libgit2 libgit2-dev build-base

RUN go get -u github.com/99designs/gqlgen github.com/vektah/gorunpkg
RUN go get -u gopkg.in/libgit2/git2go.v27
RUN go get -u github.com/gorilla/websocket github.com/rs/cors

ADD . /go/src/github.com/franxois/ggql-git/
WORKDIR /go/src/github.com/franxois/ggql-git/

# Change libgit2 version in go source
RUN sed -i 's/v26/v27/g' project/project.go

WORKDIR /go/src/github.com/franxois/ggql-git/graphql/server
RUN go build

ENTRYPOINT [ "/go/src/github.com/franxois/ggql-git/graphql/server/server" ]