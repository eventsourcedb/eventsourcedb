package eventsourcedb

import (
	"sync/atomic"
	"time"
)

const (
	OurEpoch      = 1455209242000 // OUR epoch starts from this millisecond since unix epoch
	milliSecBits  = 41
	virtShardBits = 13
	incSeqBits    = 10
	virtShards    = 1 << virtShardBits
)

var (
	autoIncCounter uint64
)

type ID uint64

func New(shardKey uint64) ID {
	now := time.Now().UnixNano() / 1000 / 1000 // milliseconds
	sinceEpoch := uint64(now - OurEpoch)
	id := sinceEpoch << (64 - milliSecBits)
	shardID := shardKey % virtShards
	id |= shardID << (64 - milliSecBits - virtShardBits)
	id |= (atomic.AddUint64(&autoIncCounter, 1) % (1 << incSeqBits))
	return ID(id)
}
