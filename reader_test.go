// Copyright 2019 kts of kettek / Ketchetwahmeegwun Tecumseh Southall. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apng

import (
	"bytes"
	"image/color"
	"os"
	"strings"
	"testing"
)

func TestIncompleteIDATOnRowBoundary(t *testing.T) {
	// The following is an invalid 1x2 grayscale PNG image. The header is OK,
	// but the zlib-compressed IDAT payload contains two bytes "\x02\x00",
	// which is only one row of data (the leading "\x02" is a row filter).
	const (
		ihdr = "\x00\x00\x00\x0dIHDR\x00\x00\x00\x01\x00\x00\x00\x02\x08\x00\x00\x00\x00\xbc\xea\xe9\xfb"
		idat = "\x00\x00\x00\x0eIDAT\x78\x9c\x62\x62\x00\x04\x00\x00\xff\xff\x00\x06\x00\x03\xfa\xd0\x59\xae"
		iend = "\x00\x00\x00\x00IEND\xae\x42\x60\x82"
	)
	_, err := Decode(strings.NewReader(pngHeader + ihdr + idat + iend))
	if err == nil {
		t.Fatal("got nil error, want non-nil")
	}
}

func TestTrailingIDATChunks(t *testing.T) {
	// The following is a valid 1x1 PNG image containing color.Gray{255} and
	// a trailing zero-length IDAT chunk (see PNG specification section 12.9):
	const (
		ihdr      = "\x00\x00\x00\x0dIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x00\x00\x00\x00\x3a\x7e\x9b\x55"
		idatWhite = "\x00\x00\x00\x0eIDAT\x78\x9c\x62\xfa\x0f\x08\x00\x00\xff\xff\x01\x05\x01\x02\x5a\xdd\x39\xcd"
		idatZero  = "\x00\x00\x00\x00IDAT\x35\xaf\x06\x1e"
		iend      = "\x00\x00\x00\x00IEND\xae\x42\x60\x82"
	)
	_, err := Decode(strings.NewReader(pngHeader + ihdr + idatWhite + idatZero + iend))
	if err != nil {
		t.Fatalf("decoding valid image: %v", err)
	}

	// Non-zero-length trailing IDAT chunks should be ignored (recoverable error).
	// The following chunk contains a single pixel with color.Gray{0}.
	const idatBlack = "\x00\x00\x00\x0eIDAT\x78\x9c\x62\x62\x00\x04\x00\x00\xff\xff\x00\x06\x00\x03\xfa\xd0\x59\xae"

	img, err := Decode(strings.NewReader(pngHeader + ihdr + idatWhite + idatBlack + iend))
	if err != nil {
		t.Fatalf("trailing IDAT not ignored: %v", err)
	}
	if img.At(0, 0) == (color.Gray{0}) {
		t.Fatal("decoded image from trailing IDAT chunk")
	}
}

func TestMultipletRNSChunks(t *testing.T) {
	/*
		The following is a valid 1x1 paletted PNG image with a 1-element palette
		containing color.NRGBA{0xff, 0x00, 0x00, 0x7f}:
			0000000: 8950 4e47 0d0a 1a0a 0000 000d 4948 4452  .PNG........IHDR
			0000010: 0000 0001 0000 0001 0803 0000 0028 cb34  .............(.4
			0000020: bb00 0000 0350 4c54 45ff 0000 19e2 0937  .....PLTE......7
			0000030: 0000 0001 7452 4e53 7f80 5cb4 cb00 0000  ....tRNS..\.....
			0000040: 0e49 4441 5478 9c62 6200 0400 00ff ff00  .IDATx.bb.......
			0000050: 0600 03fa d059 ae00 0000 0049 454e 44ae  .....Y.....IEND.
			0000060: 4260 82                                  B`.
		Dropping the tRNS chunk makes that color's alpha 0xff instead of 0x7f.
	*/
	const (
		ihdr = "\x00\x00\x00\x0dIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x03\x00\x00\x00\x28\xcb\x34\xbb"
		plte = "\x00\x00\x00\x03PLTE\xff\x00\x00\x19\xe2\x09\x37"
		trns = "\x00\x00\x00\x01tRNS\x7f\x80\x5c\xb4\xcb"
		idat = "\x00\x00\x00\x0eIDAT\x78\x9c\x62\x62\x00\x04\x00\x00\xff\xff\x00\x06\x00\x03\xfa\xd0\x59\xae"
		iend = "\x00\x00\x00\x00IEND\xae\x42\x60\x82"
	)
	for i := 0; i < 4; i++ {
		var b []byte
		b = append(b, pngHeader...)
		b = append(b, ihdr...)
		b = append(b, plte...)
		for j := 0; j < i; j++ {
			b = append(b, trns...)
		}
		b = append(b, idat...)
		b = append(b, iend...)

		var want color.Color
		m, err := Decode(bytes.NewReader(b))
		switch i {
		case 0:
			if err != nil {
				t.Errorf("%d tRNS chunks: %v", i, err)
				continue
			}
			want = color.RGBA{0xff, 0x00, 0x00, 0xff}
		case 1:
			if err != nil {
				t.Errorf("%d tRNS chunks: %v", i, err)
				continue
			}
			want = color.NRGBA{0xff, 0x00, 0x00, 0x7f}
		default:
			if err == nil {
				t.Errorf("%d tRNS chunks: got nil error, want non-nil", i)
			}
			continue
		}
		if got := m.At(0, 0); got != want {
			t.Errorf("%d tRNS chunks: got %T %v, want %T %v", i, got, got, want, want)
		}
	}
}

func TestUnknownChunkLengthUnderflow(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x06, 0xf4, 0x7c, 0x55, 0x04, 0x1a,
		0xd3, 0x11, 0x9a, 0x73, 0x00, 0x00, 0xf8, 0x1e, 0xf3, 0x2e, 0x00, 0x00,
		0x01, 0x00, 0xff, 0xff, 0xff, 0xff, 0x07, 0xf4, 0x7c, 0x55, 0x04, 0x1a,
		0xd3}
	_, err := Decode(bytes.NewReader(data))
	if err == nil {
		t.Errorf("Didn't fail reading an unknown chunk with length 0xffffffff")
	}
}

func TestGray8Transparent(t *testing.T) {
	// These bytes come from https://golang.org/issues/19553
	m, err := Decode(bytes.NewReader([]byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x0b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x85, 0x2c, 0x88,
		0x80, 0x00, 0x00, 0x00, 0x02, 0x74, 0x52, 0x4e, 0x53, 0x00, 0xff, 0x5b, 0x91, 0x22, 0xb5, 0x00,
		0x00, 0x00, 0x02, 0x62, 0x4b, 0x47, 0x44, 0x00, 0xff, 0x87, 0x8f, 0xcc, 0xbf, 0x00, 0x00, 0x00,
		0x09, 0x70, 0x48, 0x59, 0x73, 0x00, 0x00, 0x0a, 0xf0, 0x00, 0x00, 0x0a, 0xf0, 0x01, 0x42, 0xac,
		0x34, 0x98, 0x00, 0x00, 0x00, 0x07, 0x74, 0x49, 0x4d, 0x45, 0x07, 0xd5, 0x04, 0x02, 0x12, 0x11,
		0x11, 0xf7, 0x65, 0x3d, 0x8b, 0x00, 0x00, 0x00, 0x4f, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63,
		0xf8, 0xff, 0xff, 0xff, 0xb9, 0xbd, 0x70, 0xf0, 0x8c, 0x01, 0xc8, 0xaf, 0x6e, 0x99, 0x02, 0x05,
		0xd9, 0x7b, 0xc1, 0xfc, 0x6b, 0xff, 0xa1, 0xa0, 0x87, 0x30, 0xff, 0xd9, 0xde, 0xbd, 0xd5, 0x4b,
		0xf7, 0xee, 0xfd, 0x0e, 0xe3, 0xef, 0xcd, 0x06, 0x19, 0x14, 0xf5, 0x1e, 0xce, 0xef, 0x01, 0x31,
		0x92, 0xd7, 0x82, 0x41, 0x31, 0x9c, 0x3f, 0x07, 0x02, 0xee, 0xa1, 0xaa, 0xff, 0xff, 0x9f, 0xe1,
		0xd9, 0x56, 0x30, 0xf8, 0x0e, 0xe5, 0x03, 0x00, 0xa9, 0x42, 0x84, 0x3d, 0xdf, 0x8f, 0xa6, 0x8f,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	const hex = "0123456789abcdef"
	var got []byte
	bounds := m.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if r, _, _, a := m.At(x, y).RGBA(); a != 0 {
				got = append(got,
					hex[0x0f&(r>>12)],
					hex[0x0f&(r>>8)],
					' ',
				)
			} else {
				got = append(got,
					'.',
					'.',
					' ',
				)
			}
		}
		got = append(got, '\n')
	}

	const want = "" +
		".. .. .. ce bd bd bd bd bd bd bd bd bd bd e6 \n" +
		".. .. .. 7b 84 94 94 94 94 94 94 94 94 6b bd \n" +
		".. .. .. 7b d6 .. .. .. .. .. .. .. .. 8c bd \n" +
		".. .. .. 7b d6 .. .. .. .. .. .. .. .. 8c bd \n" +
		".. .. .. 7b d6 .. .. .. .. .. .. .. .. 8c bd \n" +
		"e6 bd bd 7b a5 bd bd f7 .. .. .. .. .. 8c bd \n" +
		"bd 6b 94 94 94 94 5a ef .. .. .. .. .. 8c bd \n" +
		"bd 8c .. .. .. .. 63 ad ad ad ad ad ad 73 bd \n" +
		"bd 8c .. .. .. .. 63 9c 9c 9c 9c 9c 9c 9c de \n" +
		"bd 6b 94 94 94 94 5a ef .. .. .. .. .. .. .. \n" +
		"e6 b5 b5 b5 b5 b5 b5 f7 .. .. .. .. .. .. .. \n"

	if string(got) != want {
		t.Errorf("got:\n%swant:\n%s", got, want)
	}
}

func TestDimensionOverflow(t *testing.T) {
	maxInt32AsInt := int((1 << 31) - 1)
	have32BitInts := 0 > (1 + maxInt32AsInt)

	testCases := []struct {
		src               []byte
		unsupportedConfig bool
		width             int
		height            int
	}{
		// These bytes come from https://golang.org/issues/22304
		//
		// It encodes a 2147483646 × 2147483646 (i.e. 0x7ffffffe × 0x7ffffffe)
		// NRGBA image. The (width × height) per se doesn't overflow an int64, but
		// (width × height × bytesPerPixel) will.
		{
			src: []byte{
				0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
				0x7f, 0xff, 0xff, 0xfe, 0x7f, 0xff, 0xff, 0xfe, 0x08, 0x06, 0x00, 0x00, 0x00, 0x30, 0x57, 0xb3,
				0xfd, 0x00, 0x00, 0x00, 0x15, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x62, 0x62, 0x20, 0x12, 0x8c,
				0x2a, 0xa4, 0xb3, 0x42, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x13, 0x38, 0x00, 0x15, 0x2d, 0xef,
				0x5f, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
			},
			// It's debatable whether DecodeConfig (which does not allocate a
			// pixel buffer, unlike Decode) should fail in this case. The Go
			// standard library has made its choice, and the standard library
			// has compatibility constraints.
			unsupportedConfig: true,
			width:             0x7ffffffe,
			height:            0x7ffffffe,
		},

		// The next three cases come from https://golang.org/issues/38435

		{
			src: []byte{
				0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0xb5, 0x04, 0x00, 0x00, 0xb5, 0x04, 0x08, 0x06, 0x00, 0x00, 0x00, 0xf5, 0x60, 0x2c,
				0xb8, 0x00, 0x00, 0x00, 0x15, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x62, 0x62, 0x20, 0x12, 0x8c,
				0x2a, 0xa4, 0xb3, 0x42, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x13, 0x38, 0x00, 0x15, 0x2d, 0xef,
				0x5f, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
			},
			// Here, width * height = 0x7ffea810, just under MaxInt32, but at 4
			// bytes per pixel, the number of pixels overflows an int32.
			unsupportedConfig: have32BitInts,
			width:             0x0000b504,
			height:            0x0000b504,
		},

		{
			src: []byte{
				0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
				0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x30, 0x6e, 0xc5,
				0x21, 0x00, 0x00, 0x00, 0x15, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x62, 0x62, 0x20, 0x12, 0x8c,
				0x2a, 0xa4, 0xb3, 0x42, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x13, 0x38, 0x00, 0x15, 0x2d, 0xef,
				0x5f, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
			},
			unsupportedConfig: false,
			width:             0x04000000,
			height:            0x00000001,
		},

		{
			src: []byte{
				0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
				0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0xaa, 0xd4, 0x7c,
				0xda, 0x00, 0x00, 0x00, 0x15, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x62, 0x66, 0x20, 0x12, 0x30,
				0x8d, 0x2a, 0xa4, 0xaf, 0x42, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x14, 0xd2, 0x00, 0x16, 0x00,
				0x00, 0x00,
			},
			unsupportedConfig: false,
			width:             0x08000000,
			height:            0x00000001,
		},
	}

	for i, tc := range testCases {
		cfg, err := DecodeConfig(bytes.NewReader(tc.src))
		if tc.unsupportedConfig {
			if err == nil {
				t.Errorf("i=%d: DecodeConfig: got nil error, want non-nil", i)
			} else if _, ok := err.(UnsupportedError); !ok {
				t.Fatalf("Decode: got %v (of type %T), want non-nil error (of type png.UnsupportedError)", err, err)
			}
			continue
		} else if err != nil {
			t.Errorf("i=%d: DecodeConfig: %v", i, err)
			continue
		} else if cfg.Width != tc.width {
			t.Errorf("i=%d: width: got %d, want %d", i, cfg.Width, tc.width)
			continue
		} else if cfg.Height != tc.height {
			t.Errorf("i=%d: height: got %d, want %d", i, cfg.Height, tc.height)
			continue
		}

		if nPixels := int64(cfg.Width) * int64(cfg.Height); nPixels > 0x7f000000 {
			// In theory, calling Decode would succeed, given several gigabytes
			// of memory. In practice, trying to make a []uint8 big enough to
			// hold all of the pixels can often result in OOM (out of memory).
			// OOM is unrecoverable; we can't write a test that passes when OOM
			// happens. Instead we skip the Decode call (and its tests).
			continue
		} else if testing.Short() {
			// Even for smaller image dimensions, calling Decode might allocate
			// 1 GiB or more of memory. This is usually feasible, and we want
			// to check that calling Decode doesn't panic if there's enough
			// memory, but we provide a runtime switch (testing.Short) to skip
			// these if it would OOM. See also http://golang.org/issue/5050
			// "decoding... images can cause huge memory allocations".
			continue
		}

		// Even if we don't panic, these aren't valid PNG images.
		if _, err := Decode(bytes.NewReader(tc.src)); err == nil {
			t.Errorf("i=%d: Decode: got nil error, want non-nil", i)
		}
	}

	if testing.Short() {
		t.Skip("skipping tests which allocate large pixel buffers")
	}
}

func TestReadAPNGWithDefaultFrame(t *testing.T) {
	a, err := readAPNG("tests/WithDefaultFrame.png")
	if err != nil {
		t.Error(err)
		return
	}

	if len(a.Frames) != 5 {
		t.Error("Expected 5 frames.")
		return
	}

	if !a.Frames[0].IsDefault {
		t.Error("Expected first frame to be default")
		return
	}
}

func TestReadAPNGWithoutDefaultFrame(t *testing.T) {
	a, err := readAPNG("tests/WithoutDefaultFrame.png")
	if err != nil {
		t.Error(err)
		return
	}

	if len(a.Frames) != 4 {
		t.Error("Expected 4 frames.")
		return
	}

	if a.Frames[0].IsDefault {
		t.Error("Expected first frame to not be default")
		return
	}
}

func TestReadAPNGWithMultipleIDATs(t *testing.T) {
	a, err := readAPNG("tests/MultipleIDATs.png")
	if err != nil {
		t.Error(err)
		return
	}

	if len(a.Frames) != 2 {
		t.Error("Expected 2 frames.")
		return
	}
}

func readAPNG(path string) (APNG, error) {
	f, err := os.Open(path)
	if err != nil {
		return APNG{}, err
	}
	defer f.Close()

	a, err := DecodeAll(f)
	if err != nil {
		return APNG{}, err
	}

	return a, err
}

func TestDecodeConfig(t *testing.T) {
	f, err := os.Open("tests/WithDefaultFrame.png")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	cfg, err := DecodeConfig(f)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.Width != 200 {
		t.Error("Expected 200 pixels wide.")
		return
	}
	if cfg.Height != 239 {
		t.Error("Expected 239 pixels high.")
		return
	}
	if cfg.ColorModel != color.RGBAModel {
		t.Error("Expected RGBA color model.")
		return
	}
}
