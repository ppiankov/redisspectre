package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initFlags struct {
	force bool
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate sample config",
	Long:  `Creates a sample .redisspectre.yaml with default settings.`,
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initFlags.force, "force", false, "Overwrite existing files")
	rootCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) error {
	configPath := ".redisspectre.yaml"

	if err := writeIfNotExists(configPath, sampleConfig, initFlags.force); err != nil {
		return err
	}

	fmt.Printf("Created %s\n", configPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit .redisspectre.yaml with your Redis connection details")
	fmt.Println("  2. Run: redisspectre audit")
	return nil
}

func writeIfNotExists(path, content string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists (use --force to overwrite)", path)
		}
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

const sampleConfig = `# redisspectre configuration
# See: https://github.com/ppiankov/redisspectre

# Redis connection
addr: localhost:6379
# password: ""
db: 0

# Key sampling size (number of keys to inspect)
sample_size: 10000

# Key inactivity threshold (days)
idle_days: 30

# Big key threshold (bytes, default 10MB)
big_key_size: 10485760

# Output format: text, json, sarif, spectrehub
format: text

# Audit timeout
timeout: 5m
`
