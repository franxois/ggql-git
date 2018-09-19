package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"sync"

	git "gopkg.in/libgit2/git2go.v24"
)

// git "gopkg.in/libgit2/git2go.v24" ubuntu 16.04
// git "gopkg.in/libgit2/git2go.v27" alpine 3.8
// git "gopkg.in/libgit2/git2go.v26" ubuntu 18.04

// Project contains project informations
type Project struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Path string          `json:"-"`
	Repo *git.Repository `json:"-"`
}

// GetProject get a project from a path
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

// List of projects
type List []Project

// GetProjects get a list of projects
func GetProjects(basePath string) (List, error) {
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

// CurrentBranch returns current git branch
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

// Branches return list of git branches
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

// func (p Project) GetTags() ([]string, error) {
// 	return p.Repo.Tags.ListWithMatch("v*")
// }

// GitFetch runs git fetch
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

// GitFastForward try to fast forward local develop with remote
func (p Project) GitFastForward() error {

	repo := p.Repo

	branchName, _ := p.CurrentBranch()

	fmt.Printf("%+v branch name : %s \n", p.Name, *branchName)

	if *branchName == "develop" {

		// Get remote develop
		remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + *branchName)
		if err != nil {
			return err
		}

		remoteBranchID := remoteBranch.Target()

		annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
		if err != nil {
			return err
		}

		analysis, _, err := repo.MergeAnalysis([]*git.AnnotatedCommit{annotatedCommit})
		if err != nil {
			return err
		}

		if analysis&git.MergeAnalysisUpToDate != 0 {
			fmt.Printf("%+v local develop up-to-date with remote\n", p.Name)
		} else if analysis&git.MergeAnalysisFastForward != 0 {

			branchRef, err := repo.References.Lookup("refs/heads/develop")
			if err != nil {
				fmt.Printf("Unable to find local develop branch ??%+v", err)
				return fmt.Errorf("Unable to find local develop branch ??%+v", err)
			}

			fmt.Printf("We can go fast forward %+v between from %+v to %+v \n",
				p.Name,
				branchRef.Target(),
				remoteBranchID,
			)

			// Point branch to the object
			if _, err := branchRef.SetTarget(remoteBranch.Target(), ""); err != nil {
				fmt.Printf("Error when seting target on develop %+v \n", err)
			}
		}
	}

	return nil
}

// GitUpdate runs in parallel fetch and forward
func (projects List) GitUpdate() error {

	var wg sync.WaitGroup

	projectsJobs := make(chan Project, len(projects))
	wg.Add(len(projects))

	for _, project := range projects {
		projectsJobs <- project
	}

	nbMaxConnections := 4

	for i := 0; i < nbMaxConnections; i++ {
		go func() {

			for project := range projectsJobs {
				project.GitFetch()
				project.GitFastForward()
				wg.Done()
			}
		}()
	}

	wg.Wait()
	close(projectsJobs)

	return nil
}

// Versions return list of all tag versions at semver format
func (p Project) Versions() ([]Version, error) {

	tags, err := p.Repo.Tags.ListWithMatch("v*")

	if err != nil {
		return nil, err
	}

	return tagListToVersions(tags), nil
}

// LastCandidate returns last candidate version at semver format
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

// LastRelease returns last version at semver format
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

// LastVersion returns last version at semver format
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
