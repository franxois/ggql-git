gqlgen: server/schema.graphql
	vgo generate ./server/graph/graph.go

dev: server/main.go server/graph/*.go
	make gqlgen && (cd server && vgo run main.go )
