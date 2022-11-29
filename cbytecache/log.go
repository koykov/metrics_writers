package cbytecache

import (
	"log"
	"time"
)

// LogMetrics is Log implementation of cbytecache.MetricsWriter.
//
// Don't use in production. Only for debug purposes.
type LogMetrics struct {
	key string
}

func NewLogMetrics(key string) *LogMetrics {
	m := &LogMetrics{key}
	return m
}

func (m LogMetrics) Alloc(bucket string) {
	log.Printf("cbytecache %s: alloc new arena in bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Fill(bucket string) {
	log.Printf("cbytecache %s: fill arena of bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Reset(bucket string) {
	log.Printf("cbytecache %s: reset arena of bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Release(bucket string) {
	log.Printf("cbytecache %s: release arena of bucket %s\n", m.key, bucket)
}

func (m LogMetrics) ArenaMap(bucket string, total, used, free, size uint32) {
	log.Printf("cbytecache %s: arenas mapping: total %d, used %d, free %d with size %d bytes  %s\n",
		m.key, bucket, total, used, free, size)
}

func (m LogMetrics) Set(bucket string, dur time.Duration) {
	log.Printf("cbytecache %s: set new entry to bucket %s took %s\n", m.key, bucket, dur)
}

func (m LogMetrics) Evict(bucket string) {
	log.Printf("cbytecache %s: evict entry from bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Miss(bucket string) {
	log.Printf("cbytecache %s: cache miss in bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Hit(bucket string, dur time.Duration) {
	log.Printf("cbytecache %s: cache hit in bucket %s took %s\n", m.key, bucket, dur)
}

func (m LogMetrics) Expire(bucket string) {
	log.Printf("cbytecache %s: hit expired entry in bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Corrupt(bucket string) {
	log.Printf("cbytecache %s: hit corrupted entry in bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Collision(bucket string) {
	log.Printf("cbytecache %s: keys collision in bucket %s\n", m.key, bucket)
}

func (m LogMetrics) NoSpace(bucket string) {
	log.Printf("cbytecache %s: no space in bucket %s\n", m.key, bucket)
}

func (m LogMetrics) Dump() {
	log.Printf("cbytecache %s: no space available to set new entry\n", m.key)
}

func (m LogMetrics) Load() {
	log.Printf("cbytecache %s: no space available to set new entry\n", m.key)
}

var _ = NewLogMetrics
