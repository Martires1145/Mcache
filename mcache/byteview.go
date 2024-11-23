package mcache

type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

func (b ByteView) BytesSlice() []byte {
	return cloneBytes(b.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
