package local

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

type Goroutine struct {
	data sync.Map
}

func (g *Goroutine) goId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func (g *Goroutine) Set(data any) {
	g.data.Store(g.goId(), data)
}

func (g *Goroutine) Load() (any, bool) {
	return g.data.Load(g.goId())
}
