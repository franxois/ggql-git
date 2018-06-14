gqlgen: server/schema.graphql
	go generate ./server/graph/graph.go

dev: server/main.go server/graph/*.go
	make gqlgen && (cd server && go run main.go )