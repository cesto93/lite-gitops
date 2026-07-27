package git

import (
	"errors"
	"log/slog"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Repo struct {
	repo *gogit.Repository
	path string
}

func OpenOrClone(repoURL, path, branch string) (*Repo, error) {
	r, err := gogit.PlainOpen(path)
	if err == nil {
		slog.Debug("repo already exists, opened", "path", path)
		return &Repo{repo: r, path: path}, nil
	}

	if !errors.Is(err, gogit.ErrRepositoryNotExists) {
		return nil, err
	}

	slog.Info("cloning repository", "url", repoURL, "branch", branch)
	r, err = gogit.PlainClone(path, false, &gogit.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      nil,
	})
	if err != nil {
		return nil, err
	}

	return &Repo{repo: r, path: path}, nil
}

func (r *Repo) Pull(remoteName, branch string) (bool, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return false, err
	}

	before, err := r.HeadHash()
	if err != nil {
		return false, err
	}

	err = w.Pull(&gogit.PullOptions{
		RemoteName:    remoteName,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      nil,
	})
	if err != nil {
		if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
			return false, nil
		}
		return false, err
	}

	after, _ := r.HeadHash()
	slog.Debug("pull successful", "old_hash", before, "new_hash", after)
	return true, nil
}

func (r *Repo) HeadHash() (string, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}
