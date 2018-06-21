package graph

import (
	"os"

	git "gopkg.in/libgit2/git2go.v26"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"-"`
}

func (p Project) getCurentBranch() (*string, error) {

	branchName := ""

	if _, err := os.Stat(p.Path); !os.IsNotExist(err) {
		repo, err := git.OpenRepository(p.Path)

		if err != nil {
			return nil, err
		}

		branches, _ := repo.NewBranchIterator(git.BranchLocal)

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
				branchName = name
			}

			return nil
		})

	}

	return &branchName, nil
}
