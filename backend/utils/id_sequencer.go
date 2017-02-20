package utils

import "sync/atomic"

type IdSequencer struct {
	lastId uint32
}

func NewIdSequencer() *IdSequencer {
	return &IdSequencer{
		lastId: 0,
	}
}

func (seq *IdSequencer) NextId() uint32 {
	return atomic.AddUint32(&seq.lastId, 1)
}
