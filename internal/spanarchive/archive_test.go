package spanarchive

import (
	"testing"

	"tracing/internal/spanwire"
)

type heldExporter struct{ payloads map[string][]byte }
func (e *heldExporter) Export(id string, payload []byte) { e.payloads[id] = payload }

func TestBatchesKeepIndependentPayloads(t *testing.T) {
	exporter := &heldExporter{payloads: make(map[string][]byte)}
	archive := New(&spanwire.Decoder{}, exporter)
	archive.Ingest("batch-first", "alpha")
	archive.Ingest("batch-next", "bravo")

	if got := string(exporter.payloads["batch-first"]); got != "alpha" {
		t.Errorf("first exported batch = %q, want alpha", got)
	}
	if got := string(archive.Cached("batch-first")); got != "alpha" {
		t.Errorf("first cached batch = %q, want alpha", got)
	}
	returned := archive.Cached("batch-first")
	returned[0] = 'X'
	if got := string(archive.Cached("batch-first")); got != "alpha" {
		t.Errorf("caller changed cached first batch to %q", got)
	}
}
