package tests

import (
	"goapi/conf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg, err := conf.Load()
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	assert.Equal(t, cfg.Server.Addr, "127.0.0.1:60001")
}
