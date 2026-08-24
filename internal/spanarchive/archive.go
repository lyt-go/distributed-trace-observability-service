package spanarchive

import "tracing/internal/spanwire"

type Exporter interface{ Export(string, []byte) }

type Archive struct {
	decoder *spanwire.Decoder
	exporter Exporter
	cache map[string][]byte
}

func New(decoder *spanwire.Decoder, exporter Exporter) *Archive {
	return &Archive{decoder: decoder, exporter: exporter, cache: make(map[string][]byte)}
}

func (a *Archive) Ingest(id, payload string) {
	decoded := a.decoder.Decode(payload)
	a.cache[id] = decoded
	a.exporter.Export(id, decoded)
}

func (a *Archive) Cached(id string) []byte { return a.cache[id] }
