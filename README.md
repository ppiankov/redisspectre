# redisspectre

[![CI](https://github.com/ppiankov/redisspectre/actions/workflows/ci.yml/badge.svg)](https://github.com/ppiankov/redisspectre/actions/workflows/ci.yml)

**Redis waste and hygiene auditor**

redisspectre audits Redis instances for waste and hygiene issues: memory fragmentation, idle keys, big keys, connection waste, eviction policy, persistence configuration, and slow commands. It produces actionable findings with severity levels for CI/CD gating and compliance reporting.

## What it is

- Read-only auditor for Redis instances
- Produces findings in text, JSON, SARIF, and SpectreHub formats
- Sampling-based key analysis (SCAN, never KEYS *)
- Part of the [Spectre](https://spectrehub.dev) infrastructure audit family

## What it is NOT

- Not a monitoring agent (point-in-time audit, not continuous)
- Not a data migration tool (read-only, never modifies keys)
- Not a Redis proxy or middleware
- Not a performance benchmark tool

## What it audits

| Resource | Signal | Severity |
|----------|--------|----------|
| Memory | Fragmentation ratio > 1.5 | high |
| Idle keys | Keys with no access in scan window (OBJECT IDLETIME) | medium |
| Big keys | Keys > 10MB (MEMORY USAGE) | medium |
| Connection waste | Rejected connections detected | low |
| Eviction policy | noeviction with memory near maxmemory | critical |
| Persistence | No RDB/AOF configured | high |
| Slow log | Commands > 10ms in slowlog | medium |

## Quick Start

```bash
# Generate config
redisspectre init

# Run full audit
redisspectre audit

# Audit with custom address
redisspectre audit --addr redis.example.com:6379 --password secret

# JSON output for CI/CD
redisspectre audit --format json -o report.json

# SARIF for GitHub Security tab
redisspectre audit --format sarif -o results.sarif
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

## License

MIT
