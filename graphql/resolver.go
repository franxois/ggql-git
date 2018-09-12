//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	context "context"
	"fmt"
	"os"

	project "github.com/franxois/ggql-git/project"
)

type Resolver struct{}

func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type projectResolver struct{ *Resolver }

func (r *projectResolver) CurrentBranch(ctx context.Context, obj *project.Project) (*string, error) {
	branch, err := obj.GetCurentBranch()
	return branch, err
}
func (r *projectResolver) Branches(ctx context.Context, obj *project.Project) ([]string, error) {
	return obj.GetBranches()
}
func (r *projectResolver) Versions(ctx context.Context, obj *project.Project) ([]project.Version, error) {
	return obj.GetVersions()
}

func (r *projectResolver) LastCandidate(ctx context.Context, obj *project.Project) (*project.Version, error) {
	versions, err := obj.GetVersions()

	if err != nil {
		return nil, err
	}

	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i]
		if v.IsRc {
			return &v, nil
		}
	}

	return nil, nil
}
func (r *projectResolver) LastRelease(ctx context.Context, obj *project.Project) (*project.Version, error) {
	versions, err := obj.GetVersions()

	if err != nil {
		return nil, err
	}

	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i]
		if !v.IsRc {
			return &v, nil
		}
	}

	return nil, nil
}
func (r *projectResolver) LastVersion(ctx context.Context, obj *project.Project) (*project.Version, error) {
	versions, err := obj.GetVersions()

	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, nil // no version available
	}

	return &versions[len(versions)-1], nil
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
