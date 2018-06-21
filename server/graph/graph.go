//go:generate gqlgen -typemap types.json -schema schema.graphql
package graph

import (
	context "context"
	"fmt"
	"io/ioutil"
	"os"
)

// see https://gqlgen.com/getting-started/

type MyApp struct {
	projects []Project
}

func (a *MyApp) Query_projects(ctx context.Context) ([]Project, error) {

	basePath := fmt.Sprintf("%s/src/git.fastbooking.ch/product-techno/docker-compose-attraction/projects", os.Getenv("GOPATH"))

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)

	for _, file := range files {
		path := basePath + "/" + file.Name() + "/.git"
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			// path/to/whatever does exist
			projects = append(projects, Project{ID: file.Name(), Name: file.Name(), Path: path})
		}
	}

	return projects, nil
}

func (a *MyApp) Query_project(ctx context.Context, name string) (*Project, error) {

	basePath := fmt.Sprintf("%s/src/git.fastbooking.ch/product-techno/docker-compose-attraction/projects", os.Getenv("GOPATH"))

	path := basePath + "/" + name + "/.git"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// path/to/whatever does exist
		return &Project{ID: name, Name: name, Path: path}, nil
	}

	return nil, fmt.Errorf("Project %s not found", name)
}

func (a *MyApp) Project_currentBranch(ctx context.Context, obj *Project) (*string, error) {
	branch, err := obj.getCurentBranch()
	return branch, err
}
func (a *MyApp) Project_branches(ctx context.Context, obj *Project) ([]string, error) {
	return []string{"OK??"}, nil
}
