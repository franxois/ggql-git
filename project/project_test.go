package project

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"testing"

	git "gopkg.in/libgit2/git2go.v26"
)

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
