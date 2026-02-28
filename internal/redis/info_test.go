package redis

import "testing"

func TestParseInfo(t *testing.T) {
	raw := `# Memory
used_memory:1024000
used_memory_rss:2048000
mem_fragmentation_ratio:2.00
maxmemory:0

# Clients
connected_clients:10
blocked_clients:2
`

	info := ParseInfo(raw)

	tests := []struct {
		key  string
		want string
	}{
		{"used_memory", "1024000"},
		{"used_memory_rss", "2048000"},
		{"mem_fragmentation_ratio", "2.00"},
		{"maxmemory", "0"},
		{"connected_clients", "10"},
		{"blocked_clients", "2"},
	}

	for _, tt := range tests {
		if got := info[tt.key]; got != tt.want {
			t.Errorf("ParseInfo[%q] = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestParseInfoEmpty(t *testing.T) {
	info := ParseInfo("")
	if len(info) != 0 {
		t.Errorf("ParseInfo empty string should return empty map, got %d entries", len(info))
	}
}

func TestParseInfoCommentsOnly(t *testing.T) {
	raw := `# Memory
# Server
# Stats`
	info := ParseInfo(raw)
	if len(info) != 0 {
		t.Errorf("ParseInfo comments-only should return empty map, got %d entries", len(info))
	}
}
