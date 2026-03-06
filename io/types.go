package io

import (
	"sync/atomic"
	"unsafe"

	"github.com/bitalive/chronos/mem"
)

type Iovec struct {
	Base uintptr
	Len  uintptr
}

type Msghdr struct {
	Name       uintptr
	NameLen    uint32
	_          uint32
	Iov        uintptr
	IovLen     uintptr
	Control    uintptr
	ControlLen uintptr
	Flags      int32
	_          int32
}

type Mmsghdr struct {
	MsgHdr Msghdr
	MsgLen uint32
	_      uint32
}

// Config holds configuration parameters for the Engine.
type Config struct {
	BatchSize int // Maximum number of messages in a batch
	MsgSize   int // Maximum size of a single message
}

// Stats provides nanosecond-precision monitoring without external dependencies.
// All fields are managed with atomic operations for thread-safety and zero overhead.
type Stats struct {
	RecvPackets  atomic.Uint64
	SendPackets  atomic.Uint64
	RecvBytes    atomic.Uint64
	SendBytes    atomic.Uint64
	SyscallCount atomic.Uint64
	ErrCount     atomic.Uint64
}

type Engine struct {
	Fd          int
	batch       unsafe.Pointer // Pointer to the Mmsghdr array
	payloadsPtr unsafe.Pointer // Pointer to the raw bytes (Packets)
	config      Config
	arena       *mem.Arena
	Stats       Stats
}

// GetStats returns a snapshot of current performance metrics.
func (e *Engine) GetStats() map[string]uint64 {
	return map[string]uint64{
		"recv_packets": e.Stats.RecvPackets.Load(),
		"send_packets": e.Stats.SendPackets.Load(),
		"recv_bytes":   e.Stats.RecvBytes.Load(),
		"send_bytes":   e.Stats.SendBytes.Load(),
		"syscalls":     e.Stats.SyscallCount.Load(),
		"errors":       e.Stats.ErrCount.Load(),
	}
}

// ResetStats clears all monitoring counters.
func (e *Engine) ResetStats() {
	e.Stats.RecvPackets.Store(0)
	e.Stats.SendPackets.Store(0)
	e.Stats.RecvBytes.Store(0)
	e.Stats.SendBytes.Store(0)
	e.Stats.SyscallCount.Store(0)
	e.Stats.ErrCount.Store(0)
}
