package project

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"sync"
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

	gitFetchOption := &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		}}

	var wg sync.WaitGroup

	projectsJobs := make(chan Project, len(projects))
	wg.Add(len(projects))

	for _, project := range projects {
		projectsJobs <- project
	}

	nbMaxConnections := 1

	for i := 0; i < nbMaxConnections; i++ {
		go func() {

			for project := range projectsJobs {

				fmt.Printf("%+v begin\n", project.Name)
				repo := project.Repo

				branchName, _ := project.CurrentBranch()

				fmt.Printf("%+v branch name : %s \n", project.Name, *branchName)

				if *branchName == "develop" {

					// Locate remote
					remote, err := repo.Remotes.Lookup("origin")
					if err != nil {
						t.Fatalf("%+v", err)
					}

					// Fetch changes from remote
					if err := remote.Fetch([]string{}, gitFetchOption, ""); err != nil {
						t.Fatalf("%+v", err)
					}

					fmt.Printf("%+v fetched\n", project.Name)

					// Get remote develop
					remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + *branchName)
					if err != nil {
						t.Fatalf("%+v", err)
					}

					remoteBranchID := remoteBranch.Target()

					annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
					if err != nil {
						t.Fatalf("%+v", err)
					}

					analysis, _, err := repo.MergeAnalysis([]*git.AnnotatedCommit{annotatedCommit})
					if err != nil {
						t.Fatalf("%+v", err)
					}

					fmt.Printf("%+v remote %+v\n", project.Name, analysis)

					if analysis&git.MergeAnalysisUpToDate != 0 {
						fmt.Printf("%+v remote nothing to merge ? %+v\n", project.Name, analysis&git.MergeAnalysisUpToDate)
					} else if analysis&git.MergeAnalysisFastForward != 0 {

						fmt.Printf("Going fast forward %+v\n", project.Name)

						branchRef, err := repo.References.Lookup("refs/heads/develop")
						if err != nil {
							fmt.Printf("Unable to find local develop branch ??%+v", err)
							t.Fatalf("Unable to find local develop branch ??%+v", err)
						}

						head, _ := project.Repo.Head()

						// Point branch to the object
						branchRef.SetTarget(remoteBranchID, "")
						if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
							fmt.Printf("%+v", err)
							t.Fatalf("%+v", err)
						}

					}
					// else if analysis&git.MergeAnalysisNormal != 0 {
					// 	fmt.Printf("%+v remote we could merge ? %+v\n", project.Name, analysis&git.MergeAnalysisNormal)

					// 	if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
					// 		t.Fatalf("%+v", err)
					// 	}

					// 	fmt.Printf("%+v remote we merged %+v\n", project.Name, analysis&git.MergeAnalysisNormal)

					// 	// Check for conflicts
					// 	index, err := repo.Index()
					// 	if err != nil {
					// 		t.Fatalf("%+v", err)
					// 	}
					// 	if index.HasConflicts() {
					// 		t.Fatalf("Conflicts encountered with %s. Please resolve them.", project.Name)
					// 	}
					// }
				}

				fmt.Printf("%+v exit\n", project.Name)
				wg.Done()
			}
		}()
	}

	wg.Wait()
	close(projectsJobs)
}
