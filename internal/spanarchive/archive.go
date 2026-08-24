package spanarchive

import "tracing/internal/spanwire"

// Exporter 负责将解码后的跨度内容持久化导出。
type Exporter interface{ Export(string, []byte) }

// Archive 缓存并导出解码后的跨度内容。
type Archive struct {
	decoder *spanwire.Decoder
	exporter Exporter
	cache    map[string][]byte
}

// New 创建一个 Archive，使用给定的解码器与导出器。
func New(decoder *spanwire.Decoder, exporter Exporter) *Archive {
	return &Archive{decoder: decoder, exporter: exporter, cache: make(map[string][]byte)}
}

// Ingest 解码给定载荷，按 id 缓存并导出。
func (a *Archive) Ingest(id, payload string) {
	decoded := a.decoder.Decode(payload)
	a.cache[id] = decoded
	a.exporter.Export(id, decoded)
}

// Cached 返回 id 对应缓存内容的一份副本，调用方修改返回值不会影响缓存。
func (a *Archive) Cached(id string) []byte {
	b, ok := a.cache[id]
	if !ok {
		return nil
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
