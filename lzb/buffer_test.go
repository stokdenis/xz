package lzb

import (
	"bytes"
	"testing"
)

func TestInitBuffer(t *testing.T) {
	var b buffer
	const capacity = 30
	initBuffer(&b, capacity)
	if n := b.capacity(); n != capacity {
		t.Fatalf("capacity is %d; want %d", n, capacity)
	}
	if n := b.length(); n != 0 {
		t.Fatalf("length is %d; want %d", n, 0)
	}

}

func TestNewBuffer(t *testing.T) {
	const capacity = 30
	b := newBuffer(capacity)
	if n := b.capacity(); n != capacity {
		t.Fatalf("capacity is %d; want %d", n, capacity)
	}
	if n := b.length(); n != 0 {
		t.Fatalf("length is %d; want %d", n, 0)
	}
}

func TestBuffer_Write(t *testing.T) {
	var b buffer
	const capacity = 25
	initBuffer(&b, capacity)
	p := []byte("0123456789")
	n, err := b.Write(p)
	if err != nil {
		t.Fatalf("b.Write: unexpected error %s", err)
	}
	if n != len(p) {
		t.Fatalf("b.Write returned n=%d; want %d", n, len(p))
	}
	n = b.length()
	if n != len(p) {
		t.Fatalf("b.length is %d; want %d", n, len(p))
	}
	n, err = b.Write(p)
	if err != nil {
		t.Fatalf("b.Write: unexpected error %s", err)
	}
	if n != len(p) {
		t.Fatalf("b.Write returned n=%d; want %d", n, len(p))
	}
	if n = b.length(); n != 20 {
		t.Fatalf("data length %d; want %d", n, 20)
	}
	if !bytes.Equal(b.data[:10], p) {
		t.Fatalf("first 10 byte of data wrong")
	}
	if !bytes.Equal(b.data[10:20], p) {
		t.Fatalf("second batch of 10 bytes data wrong: %q", b.data[10:])
	}
	n, err = b.Write(p)
	if err != nil {
		t.Fatalf("b.Write: unexpected error %s", err)
	}
	if n != len(p) {
		t.Fatalf("b.Write returned n=%d; want %d", n, len(p))
	}
	if b.top != 30 {
		t.Fatalf("b.top is %d; want %d", b.top, 30)
	}
	if b.bottom != 5 {
		t.Fatalf("b.bottom is %d; want %d", b.bottom, 35)
	}
	t.Logf("b.data %q", b.data)
	if !bytes.Equal(b.data[:5], p[5:]) {
		t.Fatalf("b.Write overflow problem: b.data[:5] is %q; want %q",
			b.data[:5], p[5:])
	}
	q := make([]byte, 0, 30)
	for i := 0; i < 3; i++ {
		q = append(q, p...)
	}
	n, err = b.Write(q)
	if err != nil {
		t.Fatalf("b.Write: unexpected error %s", err)
	}
	if n != len(q) {
		t.Fatalf("b.Write returned n=%d; want %d", n, len(q))
	}
	if b.top != 60 {
		t.Fatalf("b.top is %d; want %d", b.top, 60)
	}
	if !bytes.Equal(b.data[10:], q[5:20]) {
		t.Fatalf("b.data[:10] is %q; want %q", b.data[:10], q[20:])
	}
	if !bytes.Equal(b.data[:10], q[20:]) {
		t.Fatalf("b.data[:10] is %q; want %q", b.data[:10], q[20:])
	}
	n, err = b.Write([]byte{})
	if err != nil {
		t.Fatalf("b.Write: error %s", err)
	}
	if n != 0 {
		t.Fatalf("b.Write([]byte{}) returned %d; want %d", n, 0)
	}
}

func TestBuffer_Write_limit(t *testing.T) {
	b := newBuffer(20)
	b.writeLimit = 9
	p := []byte("0123456789")
	n, err := b.Write(p)
	if err != errLimit {
		t.Fatalf("b.Write error %s; want %s", err, errLimit)
	}
	if n != 9 {
		t.Fatalf("n after b.Write %d; want %d", n, 9)
	}
	b.writeLimit += 10
	n, err = b.Write(p)
	if err != nil {
		t.Fatalf("b.Write error %s", err)
	}
	if n != 10 {
		t.Fatalf("n after b.Write %d; want %d", n, 10)
	}
}

func TestBuffer_WriteByte(t *testing.T) {
	b := newBuffer(20)
	b.writeLimit = 2
	var err error
	if err = b.WriteByte(1); err != nil {
		t.Fatalf("b.WriteByte: error %s", err)
	}
	if b.top != 1 {
		t.Fatalf("after WriteByte b.top is %d; want %d", b.top, 1)
	}
	if err = b.WriteByte(1); err != nil {
		t.Fatalf("b.WriteByte: error %s", err)
	}
	if b.top != 2 {
		t.Fatalf("after WriteByte b.top is %d; want %d", b.top, 1)
	}
	if err = b.WriteByte(1); err != errLimit {
		t.Fatalf("b.WriteByte over limit error %#v; expected %#v",
			err, errLimit)
	}
}

func fillBytes(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(i)
	}
	return b
}

func TestBuffer_writeRepAt(t *testing.T) {
	b := newBuffer(10)
	b.writeLimit = 12
	p := fillBytes(5)
	var err error
	if _, err = b.Write(p); err != nil {
		t.Fatalf("Write error %s", err)
	}
	n, err := b.writeRepAt(5, 3)
	if err != nil {
		t.Fatalf("writeRepAt error %s", err)
	}
	if n != 5 {
		t.Fatalf("writeRepAt returned %d; want %d", n, 5)
	}
	w := []byte{3, 4, 3, 4, 3}
	if !bytes.Equal(b.data[5:10], w) {
		t.Fatalf("new data is %v; want %v", b.data[5:10], w)
	}
	n, err = b.writeRepAt(3, 0)
	if err != errLimit {
		t.Fatalf("b.writeRepAt returned error %v; want %v", err, errLimit)
	}
	if n != 2 {
		t.Fatalf("b.writeRepAt returned %d; want %d", n, 2)
	}
}

func TestBuffer_writeRepAt_wrap(t *testing.T) {
	b := newBuffer(5)
	p := fillBytes(7)
	var err error
	if _, err = b.Write(p); err != nil {
		t.Fatalf("Write error %s", err)
	}
	n, err := b.writeRepAt(2, 4)
	if err != nil {
		t.Fatalf("writeRepAt error %s", err)
	}
	if n != 2 {
		t.Fatalf("writeRepAt returned %d; want %d", n, 2)
	}
}

func TestBuffer_writeRepAt_errors(t *testing.T) {
	b := newBuffer(5)
	p := fillBytes(7)
	var err error
	if _, err = b.Write(p); err != nil {
		t.Fatalf("Write error %s", err)
	}
	n, err := b.writeRepAt(-2, 4)
	if err != errNegLen {
		t.Fatalf("writeRepAt error %s; want %s", err, errNegLen)
	}
	if n != 0 {
		t.Fatalf("writeRepAt returned %d; want %d", n, 0)
	}
	n, err = b.writeRepAt(1, 7)
	if err != errOffset {
		t.Fatalf("writeRepAt error %s; want %s", err, errOffset)
	}
}

func TestBuffer_equalBytes(t *testing.T) {
	b := newBuffer(10)
	if _, err := b.Write([]byte("abcabcdabcd")); err != nil {
		t.Fatalf("Write error %s", err)
	}
	tests := []struct {
		off1, off2 int64
		max, n     int
	}{
		{3, 7, 10, 4},
		{3, 7, 3, 3},
		{3, 0, 10, 0}, // index 0 is smaller then bottom
		{1, 4, 10, 2},
		{5, 9, 3, 2},
		{13, 14, 10, 0},
		{5, 14, 10, 0},
		{1, 1, 20, 10},
	}
	for _, c := range tests {
		n := b.equalBytes(c.off1, c.off2, c.max)
		if n != c.n {
			t.Errorf("b.equalBytes(%d, %d, %d) is %d; want %d",
				c.off1, c.off2, c.max, n, c.n)
		}
	}
}

func TestBuffer_setTop_panics(t *testing.T) {
	b := newBuffer(10)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("setTop: no panic on negative offset")
			}
		}()
		b.setTop(-1)
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("setTop: no panic " +
					"for writeLimit violation")
			}
		}()
		b.writeLimit = 3
		b.setTop(4)
	}()
}

func TestBuffer_index_negativeOffset(t *testing.T) {
	b := newBuffer(10)
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("index: no panic negative offset")
		}
	}()
	b.index(-1)
}