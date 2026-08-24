package spanwire

// Decoder 把跨度载荷解码为字节切片。
type Decoder struct{}

// New 返回一个新的 Decoder。
func New() *Decoder { return &Decoder{} }

// Decode 解码载荷并返回与 Decoder 内部状态完全独立的字节切片。
// 返回值在多次调用之间互不共享底层数组，调用方可安全持有与修改，
// 不会影响后续或之前的解码结果。
func (d *Decoder) Decode(payload string) []byte {
	return []byte(payload)
}
