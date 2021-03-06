//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	context "context"
	"fmt"
	"os"

	project "github.com/franxois/ggql-git/project"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Projects(ctx context.Context) ([]project.Project, error) {
	basePath := fmt.Sprintf("%s/src/git.fastbooking.ch/product-techno/docker-compose-attraction/projects", os.Getenv("GOPATH"))

	return project.GetProjects(basePath)
}
func (r *queryResolver) Project(ctx context.Context, name string) (*project.Project, error) {
	basePath := fmt.Sprintf("%s/src/git.fastbooking.ch/product-techno/docker-compose-attraction/projects", os.Getenv("GOPATH"))
	path := basePath + "/" + name + "/.git"
	return project.GetProject(name, path)
}
