package apng

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"testing"
)

// writeChunkRaw writes a single PNG chunk (length, type, data, CRC).
func writeChunkRaw(buf *bytes.Buffer, typ string, data []byte) {
	var l [4]byte
	binary.BigEndian.PutUint32(l[:], uint32(len(data)))
	buf.Write(l[:])
	c := crc32.NewIEEE()
	buf.WriteString(typ)
	c.Write([]byte(typ))
	buf.Write(data)
	c.Write(data)
	var crc [4]byte
	binary.BigEndian.PutUint32(crc[:], c.Sum32())
	buf.Write(crc[:])
}

// TestDecodeFcTLHugeDimensions verifies that an fcTL chunk advertising
// dimensions that overflow the per-frame image allocation is rejected with an
// error instead of panicking. Before the fix, DecodeAll panicked with
// "image: NewNRGBA Rectangle has huge or negative dimensions".
func TestDecodeFcTLHugeDimensions(t *testing.T) {
	var b bytes.Buffer
	b.Write([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})

	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:4], 1) // width 1
	binary.BigEndian.PutUint32(ihdr[4:8], 1) // height 1
	ihdr[8] = 8                              // bit depth
	ihdr[9] = 6                              // color type: truecolor + alpha
	writeChunkRaw(&b, "IHDR", ihdr)

	actl := make([]byte, 8)
	binary.BigEndian.PutUint32(actl[0:4], 1) // num_frames
	writeChunkRaw(&b, "acTL", actl)

	fctl := make([]byte, 26)
	binary.BigEndian.PutUint32(fctl[4:8], 0x7fffffff)  // width
	binary.BigEndian.PutUint32(fctl[8:12], 0x7fffffff) // height
	writeChunkRaw(&b, "fcTL", fctl)

	var z bytes.Buffer
	zw := zlib.NewWriter(&z)
	zw.Write(make([]byte, 8))
	zw.Close()
	writeChunkRaw(&b, "IDAT", z.Bytes())
	writeChunkRaw(&b, "IEND", nil)

	if _, err := DecodeAll(bytes.NewReader(b.Bytes())); err == nil {
		t.Fatal("expected an error for overflowing fcTL dimensions, got nil")
	}
}
