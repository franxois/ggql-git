package project

import (
	"fmt"
	"io/ioutil"
	"os"

	git "gopkg.in/libgit2/git2go.v26"
)

type Project struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Path string          `json:"-"`
	Repo *git.Repository `json:"-"`
}

func GetProject(name, path string) (*Project, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// path/to/whatever does exist

		repo, err := git.OpenRepository(path)

		if err != nil {
			return nil, err
		}

		return &Project{ID: name, Name: name, Path: path, Repo: repo}, nil
	}
	return nil, fmt.Errorf("Project not found")
}

func GetProjects(basePath string) ([]Project, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)

	for _, file := range files {
		path := basePath + "/" + file.Name() + "/.git"
		if project, err := GetProject(file.Name(), path); err == nil {
			projects = append(projects, *project)
		}
	}

	return projects, nil
}

func (p Project) GetCurentBranch() (*string, error) {

	var branchName *string

	branches, _ := p.Repo.NewBranchIterator(git.BranchLocal)

	branches.ForEach(func(b *git.Branch, t git.BranchType) error {
		name, err := b.Name()

		if err != nil {
			return err
		}

		isHead, err := b.IsHead()

		if err != nil {
			return err
		}

		if isHead {
			branchName = &name
		}

		return nil
	})

	return branchName, nil
}

func (p Project) GetBranches() ([]string, error) {

	allBranches := make([]string, 0)

	branches, _ := p.Repo.NewBranchIterator(git.BranchLocal)

	branches.ForEach(func(b *git.Branch, t git.BranchType) error {
		name, err := b.Name()

		if err != nil {
			return err
		}

		allBranches = append(allBranches, name)

		return nil
	})

	return allBranches, nil
}
