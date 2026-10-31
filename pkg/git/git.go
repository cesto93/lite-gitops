package git

import (
	"errors"
	"log/slog"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Auth holds the credentials used to access private repositories over HTTPS.
// For GitHub, Token is a personal access token (or a fine-grained token) and
// Username can be the GitHub username. When Username is empty, the
// "x-access-token" convention used by GitHub is applied.
type Auth struct {
	Username string
	Token    string
}

func (a Auth) method() transport.AuthMethod {
	if a.Token == "" {
		return nil
	}

	username := a.Username
	if username == "" {
		username = "x-access-token"
	}

	return &githttp.BasicAuth{
		Username: username,
		Password: a.Token,
	}
}

type Repo struct {
	repo *gogit.Repository
	path string
	auth Auth
}

func OpenOrClone(repoURL, path, branch string, auth Auth) (*Repo, error) {
	r, err := gogit.PlainOpen(path)
	if err == nil {
		slog.Debug("repo already exists, opened", "path", path)
		return &Repo{repo: r, path: path, auth: auth}, nil
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
		Auth:          auth.method(),
	})
	if err != nil {
		return nil, err
	}

	return &Repo{repo: r, path: path, auth: auth}, nil
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
		Auth:          r.auth.method(),
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
