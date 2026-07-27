package config

import "time"

type Config struct {
	GitRepo      string
	GitBranch    string
	GitRemote    string
	GitLocalPath string
	PollInterval time.Duration
	Verbose      bool
}
