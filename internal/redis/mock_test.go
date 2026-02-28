package redis

import (
	"context"
	"fmt"
	"time"
)

// mockClient implements RedisClient for testing.
type mockClient struct {
	infoResponses map[string]string
	scanKeys      []string
	idleTimes     map[string]time.Duration
	memoryUsages  map[string]int64
	slowLog       []SlowLogEntry
	configValues  map[string]map[string]string
	dbSize        int64
	pingErr       error
	infoErr       error
	scanErr       error
	idleTimeErr   error
	memUsageErr   error
	slowLogErr    error
	configErr     error
}

func newMockClient() *mockClient {
	return &mockClient{
		infoResponses: make(map[string]string),
		idleTimes:     make(map[string]time.Duration),
		memoryUsages:  make(map[string]int64),
		configValues:  make(map[string]map[string]string),
	}
}

func (m *mockClient) Ping(_ context.Context) error {
	return m.pingErr
}

func (m *mockClient) Info(_ context.Context, sections ...string) (string, error) {
	if m.infoErr != nil {
		return "", m.infoErr
	}
	if len(sections) > 0 {
		if resp, ok := m.infoResponses[sections[0]]; ok {
			return resp, nil
		}
	}
	return "", nil
}

func (m *mockClient) Scan(_ context.Context, cursor uint64, _ string, _ int64) ([]string, uint64, error) {
	if m.scanErr != nil {
		return nil, 0, m.scanErr
	}
	if cursor > 0 {
		return nil, 0, nil
	}
	return m.scanKeys, 0, nil
}

func (m *mockClient) ObjectIdleTime(_ context.Context, key string) (time.Duration, error) {
	if m.idleTimeErr != nil {
		return 0, m.idleTimeErr
	}
	if dur, ok := m.idleTimes[key]; ok {
		return dur, nil
	}
	return 0, fmt.Errorf("key not found: %s", key)
}

func (m *mockClient) MemoryUsage(_ context.Context, key string) (int64, error) {
	if m.memUsageErr != nil {
		return 0, m.memUsageErr
	}
	if usage, ok := m.memoryUsages[key]; ok {
		return usage, nil
	}
	return 0, fmt.Errorf("key not found: %s", key)
}

func (m *mockClient) SlowLogGet(_ context.Context, _ int64) ([]SlowLogEntry, error) {
	if m.slowLogErr != nil {
		return nil, m.slowLogErr
	}
	return m.slowLog, nil
}

func (m *mockClient) ConfigGet(_ context.Context, parameter string) (map[string]string, error) {
	if m.configErr != nil {
		return nil, m.configErr
	}
	if vals, ok := m.configValues[parameter]; ok {
		return vals, nil
	}
	return map[string]string{}, nil
}

func (m *mockClient) DBSize(_ context.Context) (int64, error) {
	return m.dbSize, nil
}

func (m *mockClient) Close() error {
	return nil
}
