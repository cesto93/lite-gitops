package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	Viper   *viper.Viper
)

var rootCmd = &cobra.Command{
	Use:   "lite-gitops",
	Short: "A lightweight GitOps agent",
	Long: `Lite GitOps continuously syncs configuration from a Git repository
and applies it to the local system.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default ~/.lite-gitops/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
}

func initConfig(flags *pflag.FlagSet) error {
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		v.AddConfigPath(home + "/.lite-gitops")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetDefault("git.repository", "")
	v.SetDefault("git.branch", "main")
	v.SetDefault("git.remote", "origin")
	v.SetDefault("git.local_path", "/var/lib/lite-gitops/state")
	v.SetDefault("compose.files", []string{"docker-compose.yml"})
	v.SetDefault("compose.project_name", "")
	v.SetDefault("compose.build", false)
	v.SetDefault("poll_interval", "30s")
	v.SetDefault("verbose", false)

	flags.Visit(func(f *pflag.Flag) {
		v.BindPFlag(f.Name, flags.Lookup(f.Name))
	})

	if v.GetBool("verbose") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	Viper = v
	return nil
}
