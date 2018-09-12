# Projects manager with Go, React and GraphQL

## install

### Without go dep

```
go get -u github.com/99designs/gqlgen github.com/vektah/gorunpkg github.com/hashicorp/golang-lru
go get -u github.com/oxequa/realize
cd graphql && gqlgen -v

#To run
realize start
```
### With go dep

can't make it work ...

```
rm $GOPATH/bin/gqlgen
rm -rf $GOPATH/src/github.com/99designs/gqlgen
go get -u github.com/golang/dep/cmd/dep
dep init
dep ensure
go generate ./...
```

## front

<https://github.com/wmonk/create-react-app-typescript>

### Reason-apollo

<https://github.com/apollographql/reason-apollo>

send-introspection-query <http://127.0.0.1:8080/query>
