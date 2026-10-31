package config

import "time"

type ComposeConfig struct {
	Files       []string
	ProjectName string
	Build       bool
}

type Config struct {
	GitRepo      string
	GitBranch    string
	GitRemote    string
	GitLocalPath string
	GitUsername  string
	GitToken     string
	PollInterval time.Duration
	Verbose      bool
	Compose      ComposeConfig
}
