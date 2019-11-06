package deploy

import (
	"io"
	"time"
)

type bufferingWriter struct {
	updateFunc func([]byte)
	lastUpdate time.Time
	bufData    []byte
}

var _ io.Writer = &bufferingWriter{}

func newBufferingWriter(updateFunc func([]byte)) *bufferingWriter {
	return &bufferingWriter{updateFunc: updateFunc}
}

func (w *bufferingWriter) Write(data []byte) (int, error) {
	w.bufData = append(w.bufData, data...)

	now := time.Now()
	if now.Sub(w.lastUpdate) > 5*time.Second {
		w.lastUpdate = now
		w.updateFunc(w.bufData)
		w.bufData = nil
	}

	return len(data), nil
}
