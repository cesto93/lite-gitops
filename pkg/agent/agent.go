package agent

import (
	"context"
	"log/slog"
	"time"

	"lite-gitops/pkg/config"
	"lite-gitops/pkg/git"
)

type Agent struct {
	cfg      *config.Config
	repo     *git.Repo
	lastHash string
}

func New(cfg *config.Config) (*Agent, error) {
	return &Agent{cfg: cfg}, nil
}

func (a *Agent) Run(ctx context.Context) error {
	slog.Info("agent starting",
		"repo", a.cfg.GitRepo,
		"branch", a.cfg.GitBranch,
		"poll_interval", a.cfg.PollInterval,
	)

	if err := a.ensureRepo(); err != nil {
		return err
	}

	ticker := time.NewTicker(a.cfg.PollInterval)
	defer ticker.Stop()

	for {
		hash, err := a.repo.HeadHash()
		if err != nil {
			slog.Error("failed to get HEAD hash", "error", err)
		} else {
			slog.Debug("current commit", "hash", hash)
			if a.lastHash == "" {
				a.lastHash = hash
				a.reconcile(hash)
			}
		}

		select {
		case <-ticker.C:
			a.sync()
		case <-ctx.Done():
			slog.Info("agent shutting down")
			return nil
		}
	}
}

func (a *Agent) ensureRepo() error {
	r, err := git.OpenOrClone(a.cfg.GitRepo, a.cfg.GitLocalPath, a.cfg.GitBranch)
	if err != nil {
		return err
	}
	a.repo = r
	return nil
}

func (a *Agent) sync() {
	changed, err := a.repo.Pull(a.cfg.GitRemote, a.cfg.GitBranch)
	if err != nil {
		slog.Error("pull failed", "error", err)
		return
	}

	if !changed {
		slog.Debug("no new commits")
		return
	}

	hash, err := a.repo.HeadHash()
	if err != nil {
		slog.Error("failed to get HEAD hash after pull", "error", err)
		return
	}

	slog.Info("new commit detected", "old_hash", a.lastHash, "new_hash", hash)
	a.lastHash = hash
	a.reconcile(hash)
}

func (a *Agent) reconcile(hash string) {
	slog.Info("reconciling state", "hash", hash)
}
