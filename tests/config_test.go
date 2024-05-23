package tests

import (
	"goapi/conf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, conf.Get().Server.Addr, "127.0.0.1:60001")
}
