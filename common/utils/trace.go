package utils

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"sync"
)

func NewTraceId() string {
	var seed int64
	var s sync.Mutex
	s.Lock()
	defer s.Unlock()
	_ = binary.Read(crand.Reader, binary.LittleEndian, &seed)
	tid := [16]byte{}
	rand.New(rand.NewSource(seed)).Read(tid[:])
	return hex.EncodeToString(tid[:])
}
