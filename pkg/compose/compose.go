package compose

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"

	"lite-gitops/pkg/config"
)

type Manager struct {
	cfg    config.ComposeConfig
	repoPath string
}

func NewManager(cfg config.ComposeConfig, repoPath string) *Manager {
	return &Manager{cfg: cfg, repoPath: repoPath}
}

func (m *Manager) Deploy() error {
	slog.Info("deploying docker compose", "project", m.cfg.ProjectName)

	args := []string{"compose"}

	for _, f := range m.cfg.Files {
		args = append(args, "-f", filepath.Join(m.repoPath, f))
	}

	if m.cfg.ProjectName != "" {
		args = append(args, "-p", m.cfg.ProjectName)
	}

	pullArgs := append(args, "pull")
	if err := runCommand("docker", pullArgs); err != nil {
		return fmt.Errorf("docker compose pull failed: %w", err)
	}

	upArgs := append(args, "up", "-d")
	if m.cfg.Build {
		upArgs = append(upArgs, "--build")
	}
	if err := runCommand("docker", upArgs); err != nil {
		return fmt.Errorf("docker compose up failed: %w", err)
	}

	slog.Info("docker compose deployed successfully", "project", m.cfg.ProjectName)
	return nil
}

func runCommand(name string, args []string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	slog.Debug("running command", "cmd", name, "args", args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w\nstderr: %s", name, err, stderr.String())
	}

	if stdout.Len() > 0 {
		slog.Debug("command output", "output", stdout.String())
	}
	return nil
}
