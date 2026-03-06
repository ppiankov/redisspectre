## Install

```bash
# Homebrew
brew install ppiankov/tap/redisspectre

# Go
go install github.com/ppiankov/redisspectre/cmd/redisspectre@latest

# Binary: download from GitHub Releases
# https://github.com/ppiankov/redisspectre/releases
```


## Usage

### Commands

```bash
redisspectre audit    # Full Redis audit
redisspectre init     # Generate sample config
redisspectre version  # Print version information
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--addr` | localhost:6379 | Redis address (host:port) |
| `--password` | (empty) | Redis password (or REDIS_PASSWORD env) |
| `--db` | 0 | Redis database number |
| `--format` | text | Output format: text, json, sarif, spectrehub |
| `-o, --output` | stdout | Output file path |
| `--sample-size` | 10000 | Number of keys to sample |
| `--idle-days` | 30 | Key inactivity threshold (days) |
| `--big-key-size` | 10485760 | Big key threshold (bytes) |
| `--timeout` | 5m | Audit timeout |
| `-v, --verbose` | false | Enable verbose logging |

### Configuration

Create `.redisspectre.yaml` (or run `redisspectre init`):

```yaml
addr: localhost:6379
db: 0
sample_size: 10000
idle_days: 30
big_key_size: 10485760
format: text
timeout: 5m
```


## Architecture

- **Single binary** — no dependencies, no server-side components
- **Read-only** — uses INFO, SCAN, OBJECT, MEMORY, SLOWLOG, CONFIG GET
- **Sampling-based** — never runs KEYS *, uses SCAN with count limits
- **Concurrent** — parallel auditors with bounded concurrency


## Project Status

**Status: Alpha** | v0.1.0

