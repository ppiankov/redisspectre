package commands

import (
	"log/slog"
	"os"

	"github.com/ppiankov/redisspectre/internal/config"
	"github.com/ppiankov/redisspectre/internal/logging"
	"github.com/spf13/cobra"
)

var (
	verbose  bool
	addr     string
	password string
	db       int
	version  string
	commit   string
	date     string
	cfg      config.Config
)

var rootCmd = &cobra.Command{
	Use:   "redisspectre",
	Short: "redisspectre — Redis waste and hygiene auditor",
	Long: `redisspectre audits Redis instances for waste and hygiene issues: memory
fragmentation, idle keys, big keys, connection waste, eviction policy,
persistence configuration, and slow commands.

Read-only: never modifies data. Uses INFO, SCAN, OBJECT IDLETIME,
MEMORY USAGE, SLOWLOG GET, and CONFIG GET.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logging.Init(verbose)
		loaded, err := config.Load(".")
		if err != nil {
			slog.Warn("Failed to load config file", "error", err)
		} else {
			cfg = loaded
		}
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command with injected build info.
func Execute(v, c, d string) error {
	version = v
	commit = c
	date = d
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "localhost:6379", "Redis address (host:port)")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Redis password (or REDIS_PASSWORD env)")
	rootCmd.PersistentFlags().IntVar(&db, "db", 0, "Redis database number")
	rootCmd.AddCommand(versionCmd)
}

func resolvePassword() string {
	if password != "" {
		return password
	}
	if envPwd := os.Getenv("REDIS_PASSWORD"); envPwd != "" {
		return envPwd
	}
	return cfg.Password
}

func resolveAddr() string {
	if addr != "localhost:6379" {
		return addr
	}
	if cfg.Addr != "" {
		return cfg.Addr
	}
	return addr
}
