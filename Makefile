gqlgen: ./server/graph/schema.graphql
	vgo generate ./server/graph/graph.go

dev:
	realize start
