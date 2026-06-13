package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// NewTraceId 生成 16 字节随机十六进制 trace id。
// 直接用 crypto/rand，避免 math/rand 的可预测性与多余的同步开销。
func NewTraceId() string {
	var buf [16]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}
