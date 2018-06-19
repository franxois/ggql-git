gqlgen: server/schema.graphql
	vgo generate ./server/graph/graph.go

dev:
	realize start
