# redisspectre

[![CI](https://github.com/ppiankov/redisspectre/actions/workflows/ci.yml/badge.svg)](https://github.com/ppiankov/redisspectre/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ppiankov/redisspectre)](https://goreportcard.com/report/github.com/ppiankov/redisspectre)
[![ANCC](https://img.shields.io/badge/ANCC-compliant-brightgreen)](https://ancc.dev)

**redisspectre** — Redis waste and hygiene auditor. Part of [SpectreHub](https://github.com/ppiankov/spectrehub).

## What it is

- Audits Redis instances for memory fragmentation, idle keys, big keys, and connection waste
- Checks eviction policy, persistence configuration, and slow commands
- Uses sampling-based key analysis (SCAN, never KEYS *)
- Each finding includes severity for CI/CD gating
- Outputs text, JSON, SARIF, and SpectreHub formats

## What it is NOT

- Not a monitoring agent — point-in-time audit
- Not a data migration tool — read-only, never modifies keys
- Not a Redis proxy or middleware
- Not a performance benchmark tool

## Quick start

### Homebrew

```sh
brew tap ppiankov/tap
brew install redisspectre
```

### From source

```sh
git clone https://github.com/ppiankov/redisspectre.git
cd redisspectre
make build
```

### Usage

```sh
redisspectre audit --addr localhost:6379 --format json
```

## CLI commands

| Command | Description |
|---------|-------------|
| `redisspectre audit` | Audit Redis instance for waste and hygiene issues |
| `redisspectre version` | Print version |

## SpectreHub integration

redisspectre feeds Redis hygiene findings into [SpectreHub](https://github.com/ppiankov/spectrehub) for unified visibility across your infrastructure.

```sh
spectrehub collect --tool redisspectre
```

## Safety

redisspectre operates in **read-only mode**. It inspects and reports — never modifies, deletes, or alters your keys.

## Documentation

| Document | Contents |
|----------|----------|
| [CLI Reference](docs/cli-reference.md) | Full command reference, flags, and configuration |

## License

MIT — see [LICENSE](LICENSE).

---

Built by [Obsta Labs](https://obstalabs.dev)
