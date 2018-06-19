//go:generate gqlgen -typemap types.json -schema schema.graphql
package graph

import (
	context "context"
	"fmt"
	"io/ioutil"
	"os"
)

// see https://gqlgen.com/getting-started/

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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
		if _, err := os.Stat(basePath + "/" + file.Name() + "/.git"); !os.IsNotExist(err) {
			// path/to/whatever does exist
			projects = append(projects, Project{ID: file.Name(), Name: file.Name()})
		}
	}

	return projects, nil
}

func (a *MyApp) Project_currentBranch(ctx context.Context, obj *Project) (*string, error) {
	s := "Ok chief"
	return &s, nil
}
func (a *MyApp) Project_branches(ctx context.Context, obj *Project) ([]string, error) {
	return []string{"OK??"}, nil
}
