package cmd

import (
	"context"
	"log/slog"
	"time"

	"lite-gitops/pkg/agent"
	"lite-gitops/pkg/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the GitOps agent loop",
	Long: `Clone (or pull) the configured Git repository on a polling interval
and reconcile any configuration changes.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := &config.Config{
			GitRepo:      Viper.GetString("git.repository"),
			GitBranch:    Viper.GetString("git.branch"),
			GitRemote:    Viper.GetString("git.remote"),
			GitLocalPath: Viper.GetString("git.local_path"),
			Verbose:      Viper.GetBool("verbose"),
			Compose: config.ComposeConfig{
				Files:       Viper.GetStringSlice("compose.files"),
				ProjectName: Viper.GetString("compose.project_name"),
				Build:       Viper.GetBool("compose.build"),
			},
		}

		interval, err := time.ParseDuration(Viper.GetString("poll_interval"))
		if err != nil {
			return err
		}
		cfg.PollInterval = interval

		if cfg.GitRepo == "" {
			slog.Error("git.repository is required")
			return nil
		}

		ag, err := agent.New(cfg)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return ag.Run(ctx)
	},
}
