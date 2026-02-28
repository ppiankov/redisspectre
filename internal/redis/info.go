package redis

import (
	"strings"
)

// ParseInfo parses Redis INFO output into a key-value map.
// INFO output has lines like "key:value" and section headers like "# Memory".
func ParseInfo(raw string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
