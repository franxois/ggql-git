//go:generate gqlgen

package graph

import (
	context "context"
	"fmt"
	"os"
)

// see https://gqlgen.com/getting-started/

type MyApp struct {
	projects []Project
}

var basePath string

func init() {
	basePath = fmt.Sprintf("%s/src/git.fastbooking.ch/product-techno/docker-compose-attraction/projects", os.Getenv("GOPATH"))
}

func (a *MyApp) Query_projects(ctx context.Context) ([]Project, error) {
	return getProjects(basePath)
}

func (a *MyApp) Query_project(ctx context.Context, name string) (*Project, error) {
	path := basePath + "/" + name + "/.git"
	return getProject(name, path)
}

func (a *MyApp) Project_currentBranch(ctx context.Context, obj *Project) (*string, error) {
	branch, err := obj.getCurentBranch()
	return branch, err
}
func (a *MyApp) Project_branches(ctx context.Context, obj *Project) ([]*string, error) {
	return obj.getBranches()
}
