package statsd

import (
	"bufio"
	"bytes"
	"github.com/bmizerany/assert"
	"io"
	"testing"
	"time"
)

type ReadWriter struct {
	io.Reader
	io.Writer
}

func fake() io.ReadWriter {
	buffer := bytes.NewBufferString("")
	return &ReadWriter{buffer, buffer}
}

func readData(rw *bufio.ReadWriter) string {
	data, _, _ := rw.ReadLine()
	return string(data)
}

func TestIncrement(t *testing.T) {
	c := newClient("<fake>", fake())
	err := c.Increment("incr", 1, 1)
	assert.Equal(t, err, nil)
	data := readData(c.rw)
	assert.Equal(t, data, "incr:1|c")
}

func TestDecrement(t *testing.T) {
	c := newClient("<fake>", fake())
	err := c.Decrement("decr", 1, 1)
	assert.Equal(t, err, nil)
	data := readData(c.rw)
	assert.Equal(t, data, "decr:-1|c")
}

func TestIncrementRate(t *testing.T) {
	c := newClient("<fake>", fake())
	err := c.Increment("incr", 1, 0.99)
	assert.Equal(t, err, nil)
	data := readData(c.rw)
	assert.Equal(t, data, "incr:1|c|@0.99")
}

func TestGauge(t *testing.T) {
	c := newClient("<fake>", fake())
	err := c.Gauge("gauge", 300, 1)
	assert.Equal(t, err, nil)
	data := readData(c.rw)
	assert.Equal(t, data, "gauge:300|g")
}

func TestMilliseconds(t *testing.T) {
	msec, _ := time.ParseDuration("350ms")
	assert.Equal(t, 350, millisecond(msec))
	sec, _ := time.ParseDuration("5s")
	assert.Equal(t, 5000, millisecond(sec))
	nsec, _ := time.ParseDuration("50ns")
	assert.Equal(t, 0, millisecond(nsec))
}

func TestTiming(t *testing.T) {
	c := newClient("<fake>", fake())
	err := c.Timing("timing", 350, 1)
	assert.Equal(t, err, nil)
	data := readData(c.rw)
	assert.Equal(t, data, "timing:350|ms")
}

func TestTime(t *testing.T) {
	c := newClient("<fake>", fake())
	err := c.Time("time", 1, func() { time.Sleep(50e6) })
	assert.Equal(t, err, nil)
}
