package spanwire

type Decoder struct{ scratch []byte }

func (d *Decoder) Decode(payload string) []byte {
	d.scratch = append(d.scratch[:0], payload...)
	return d.scratch
}
