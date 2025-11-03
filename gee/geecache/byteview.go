package geecache

type ByteView struct {
	b []byte
}

func NewByteView(s string) ByteView {
	return ByteView{b: []byte(s)}
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return  string(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

