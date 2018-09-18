package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	git "gopkg.in/libgit2/git2go.v24"
)

// git "gopkg.in/libgit2/git2go.v24" ubuntu 16.04
// git "gopkg.in/libgit2/git2go.v27" alpine 3.8
// git "gopkg.in/libgit2/git2go.v26" ubuntu 18.04

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

func (p Project) CurrentBranch() (*string, error) {

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

func (p Project) Branches() ([]string, error) {

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

func (p Project) GetTags() ([]string, error) {
	return p.Repo.Tags.ListWithMatch("v*")
}

func (p Project) GitFetch() error {

	repo := p.Repo

	// Locate remote
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	gitFetchOption := &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, gitFetchOption, ""); err != nil {
		return err
	}

	return nil
}

func (p Project) Versions() ([]Version, error) {

	tags, err := p.Repo.Tags.ListWithMatch("v*")

	if err != nil {
		return nil, err
	}

	return tagListToVersions(tags), nil
}

func (p Project) LastCandidate() (*Version, error) {
	versions, err := p.Versions()

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
func (p Project) LastRelease() (*Version, error) {
	versions, err := p.Versions()

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
func (p Project) LastVersion() (*Version, error) {
	versions, err := p.Versions()

	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, nil // no version available
	}

	return &versions[len(versions)-1], nil
}

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	ret, cred := git.NewCredSshKey("git", usr.HomeDir+"/.ssh/id_rsa.pub", usr.HomeDir+"/.ssh/id_rsa", "")
	return git.ErrorCode(ret), &cred
}

// Made this one just return 0 during troubleshooting...
func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}
