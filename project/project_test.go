package project

import (
	"fmt"
	"os"
	"testing"

	git "gopkg.in/libgit2/git2go.v24"
)

func TestPlayWithRepo(t *testing.T) {

	cloneOptions := &git.CloneOptions{}
	// use FetchOptions instead of directly RemoteCallbacks
	// https://github.com/libgit2/git2go/commit/36e0a256fe79f87447bb730fda53e5cbc90eb47c
	cloneOptions.FetchOptions = &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}
	if _, err := os.Stat("tmp"); os.IsNotExist(err) && false {
		repo, err := git.Clone("git@github.com:franxois/vuex-modules.git", "tmp", cloneOptions)
		if err != nil {
			panic(err)
		}

		s, _ := repo.StatusFile(".")

		fmt.Printf("Ho hi ? %+v\n", s)
	}
}

func TestPlayWithLocalRepo(t *testing.T) {

	// cloneOptions := &git.CloneOptions{}
	// cloneOptions.FetchOptions = &git.FetchOptions{
	// 	RemoteCallbacks: git.RemoteCallbacks{
	// 		CredentialsCallback:      credentialsCallback,
	// 		CertificateCheckCallback: certificateCheckCallback,
	// 	},
	// }

	basePath := fmt.Sprintf("%s/src/git.fastbooking.ch/product-techno/docker-compose-attraction/projects", os.Getenv("GOPATH"))

	projects, _ := GetProjects(basePath)

	projects.GitUpdate()
}
