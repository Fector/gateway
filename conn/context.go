package conn

import (
	"sync/atomic"
	"time"
)

type Context struct {
	SequenceNumber *uint32
	LastPacketTime *int64
}

func (c *Context) NextSequence() uint32 {
	return atomic.AddUint32(c.SequenceNumber, 1)
}

func (c *Context) UpdateTime() {
	atomic.StoreInt64(c.LastPacketTime, time.Now().UnixNano())
}
