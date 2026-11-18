// Package perfmark emits opt-in, request-local timing marks used to decompose
// create / pause / resume latency into per-stage intervals.
//
// It is disabled unless KUASAR_PERF=1. Every mark is one line on stderr, which
// systemd/journald captures with its own timestamp:
//
//	KUASAR_PERF {"comp":"conductor","op":"pause","rid":"...","sid":"...","evt":"api.recv","off_ms":0.123,"wall_ns":1763...}
//
// off_ms is the offset from the mark set's origin (monotonic, best for
// intervals); wall_ns is CLOCK_REALTIME, used only to align marks that come
// from different processes. Adjacent offsets delimit intervals; cumulative
// offsets must never be added together.
package perfmark

import (
	"encoding/json"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var enabled = os.Getenv("KUASAR_PERF") == "1"
var ridSeq atomic.Uint64

// Enabled reports whether marks are being emitted.
func Enabled() bool { return enabled }

// Marks is one operation's mark sequence. The zero value is usable and does
// nothing.
type Marks struct {
	comp   string
	op     string
	rid    string
	origin time.Time
	sid    string
	run    string
}

// Start begins a mark set for one request/operation. comp is the component
// ("conductor", "runner", "runtime"), op the operation ("create", "resume",
// "pause").
func Start(comp, op string) *Marks {
	if !enabled {
		return nil
	}
	return &Marks{
		comp:   comp,
		op:     op,
		rid:    comp + "-" + strconv.FormatUint(ridSeq.Add(1), 10),
		origin: time.Now(),
	}
}

// ID returns the mark set's request id ("" when disabled).
func (m *Marks) ID() string {
	if m == nil {
		return ""
	}
	return m.rid
}

// SID tags subsequent marks with the sandbox id.
func (m *Marks) SID(sid string) *Marks {
	if m == nil {
		return nil
	}
	m.sid = sid
	return m
}

// Run tags subsequent marks with the runner run id.
func (m *Marks) Run(run string) *Marks {
	if m == nil {
		return nil
	}
	m.run = run
	return m
}

// Mark records a cumulative checkpoint at the current instant.
func (m *Marks) Mark(evt string) {
	if m == nil {
		return
	}
	m.emit(evt, time.Since(m.origin), "cumulative")
}

// Dur records an interval that ended now and started at begin.
func (m *Marks) Dur(evt string, begin time.Time) {
	if m == nil {
		return
	}
	m.emit(evt, time.Since(begin), "interval")
}

// MarkAt records a cumulative checkpoint for a timestamp already taken.
func (m *Marks) MarkAt(evt string, at time.Time) {
	if m == nil {
		return
	}
	m.emit(evt, at.Sub(m.origin), "cumulative")
}

// MarkKV records a cumulative checkpoint together with extra numeric fields;
// the analyzer reads them as annotations, never as intervals.
func (m *Marks) MarkKV(evt string, kv map[string]any) {
	if m == nil {
		return
	}
	rec := m.base(evt, time.Since(m.origin), "cumulative")
	for k, v := range kv {
		rec[k] = v
	}
	write(rec)
}

func (m *Marks) emit(evt string, off time.Duration, kind string) {
	write(m.base(evt, off, kind))
}

func (m *Marks) base(evt string, off time.Duration, kind string) map[string]any {
	rec := map[string]any{
		"pid":     os.Getpid(),
		"comp":    m.comp,
		"op":      m.op,
		"rid":     m.rid,
		"evt":     evt,
		"kind":    kind,
		"off_ms":  float64(off) / 1e6,
		"wall_ns": time.Now().UnixNano(),
	}
	if m.sid != "" {
		rec["sid"] = m.sid
	}
	if m.run != "" {
		rec["run"] = m.run
	}
	return rec
}

func write(rec map[string]any) {
	body, err := json.Marshal(rec)
	if err != nil {
		return
	}
	_, _ = os.Stderr.WriteString("KUASAR_PERF " + string(body) + "\n")
}
