package hash

import (
	"unsafe"

	"github.com/bitalive/chronos/hash"
)

type Slot struct {
	Hash  uint64
	Value uint64
}

type Table struct {
	Slots []Slot
	Mask  uint64
}

func NewTable(buffer []byte) *Table {
	size := len(buffer) / int(unsafe.Sizeof(Slot{}))
	return &Table{
		Slots: unsafe.Slice((*Slot)(unsafe.Pointer(&buffer[0])), size),
		Mask:  uint64(size - 1),
	}
}

func (t *Table) Put(key []byte, val uint64) {
	h64 := hash.WyHash(key, 0)
	idx := h64 & t.Mask
	for {
		if t.Slots[idx].Hash == 0 {
			t.Slots[idx].Hash = h64
			t.Slots[idx].Value = val
			return
		}
		idx = (idx + 1) & t.Mask
	}
}

func (t *Table) Get(key []byte) uint64 {
	h64 := hash.WyHash(key, 0)
	idx := h64 & t.Mask
	for {
		s := t.Slots[idx]
		if s.Hash == h64 {
			return s.Value
		}
		if s.Hash == 0 {
			return 0
		}
		idx = (idx + 1) & t.Mask
	}
}
