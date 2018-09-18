package project

import (
	"fmt"
	"os"
	"sync"
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

				project.GitFetch()

				repo := project.Repo

				branchName, _ := project.CurrentBranch()

				fmt.Printf("%+v branch name : %s \n", project.Name, *branchName)

				if *branchName == "develop" {

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

						branchRef, err := repo.References.Lookup("refs/heads/develop")
						if err != nil {
							fmt.Printf("Unable to find local develop branch ??%+v", err)
							t.Fatalf("Unable to find local develop branch ??%+v", err)
						}

						head, _ := project.Repo.Head()

						fmt.Printf("We can go fast forward %+v between %+v , %+v and %+v\n",
							project.Name,
							branchRef.Target(),
							remoteBranchID,
							head.Target())

						// Point branch to the object
						if _, err := branchRef.SetTarget(remoteBranch.Target(), ""); err != nil {
							fmt.Printf("Error when seting target on develop %+v \n", err)
						}
						if _, err := head.SetTarget(remoteBranch.Target(), ""); err != nil {
							fmt.Printf("Error when seting target on head %+v \n", err)
						}
					}
				}

				fmt.Printf("%+v exit\n", project.Name)
				wg.Done()
			}
		}()
	}

	wg.Wait()
	close(projectsJobs)
}
