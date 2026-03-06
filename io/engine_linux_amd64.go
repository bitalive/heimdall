//go:build linux && amd64
// +build linux,amd64

package io

import (
	"syscall"
	"unsafe"

	"github.com/bitalive/chronos/mem"
)

// Internal Assembly calls
func Heimdall_Recv(fd uintptr, msgbatch uintptr, vlen int) (int, int)
func Heimdall_Send(fd uintptr, msgbatch uintptr, vlen int) (int, int)

// NewEngine creates a new Heimdall Network Engine for Linux AMD64
func NewEngine(port int, cfg Config) (*Engine, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}

	arena := mem.NewArena(16 * 1024 * 1024)
	batchPtr, payloads := PrepareBatch(arena, cfg.BatchSize, cfg.MsgSize)

	return &Engine{
		Fd:          fd,
		batch:       unsafe.Pointer(batchPtr),
		payloadsPtr: unsafe.Pointer(unsafe.SliceData(payloads)),
		config:      cfg,
		arena:       arena,
	}, nil
}

// Poll fetches the next batch of network packets.
func (e *Engine) Poll() int {
	e.Stats.SyscallCount.Add(1)

	// Call Assembly Layer
	n, err := Heimdall_Recv(uintptr(e.Fd), uintptr(e.batch), e.config.BatchSize)

	// Check for anomalies from Kernel (Safety Guard)
	if n > e.config.BatchSize {
		n = e.config.BatchSize
	}

	if n > 0 {
		e.Stats.RecvPackets.Add(uint64(n))
		e.Stats.RecvBytes.Add(uint64(n * e.config.MsgSize))
	} else if err != 0 && err != 11 { // 11 is EAGAIN
		e.Stats.ErrCount.Add(1)
	}

	return n
}

// Send dispatches n frames back to clients using high-perf Assembly.
func (e *Engine) Send(n int) int {
	e.Stats.SyscallCount.Add(1)
	m, err := Heimdall_Send(uintptr(e.Fd), uintptr(e.batch), n)
	if m > 0 {
		e.Stats.SendPackets.Add(uint64(m))
		e.Stats.SendBytes.Add(uint64(m * e.config.MsgSize))
	} else if err != 0 {
		e.Stats.ErrCount.Add(1)
	}
	return m
}

// GetFrame returns a raw byte slice for a specific index in the current batch.
func (e *Engine) GetFrame(index int) []byte {
	if index >= e.config.BatchSize {
		return nil
	}
	ptr := unsafe.Add(unsafe.Pointer(e.payloadsPtr), index*e.config.MsgSize)
	return unsafe.Slice((*byte)(ptr), e.config.MsgSize)
}

// PrepareBatch sets up a mmsghdr array in the Arena
func PrepareBatch(arena *mem.Arena, count int, msgSize int) (unsafe.Pointer, []byte) {
	mmsgSize := int(unsafe.Sizeof(Mmsghdr{}))
	mmsgData := arena.Alloc(count * mmsgSize)
	payloads := arena.Alloc(count * msgSize)

	mmsgs := unsafe.Slice((*Mmsghdr)(unsafe.Pointer(unsafe.SliceData(mmsgData))), count)
	for i := 0; i < count; i++ {
		iovBuf := arena.Alloc(int(unsafe.Sizeof(Iovec{})))
		iov := (*Iovec)(unsafe.Pointer(unsafe.SliceData(iovBuf)))
		iov.Base = uintptr(unsafe.Pointer(unsafe.SliceData(payloads[i*msgSize:])))
		iov.Len = uintptr(msgSize)
		mmsgs[i].MsgHdr.Iov = uintptr(unsafe.Pointer(iov))
		mmsgs[i].MsgHdr.IovLen = 1
	}
	return unsafe.Pointer(unsafe.SliceData(mmsgData)), payloads
}
