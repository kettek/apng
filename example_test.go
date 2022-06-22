package apng_test

import (
	"encoding/base64"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/kettek/apng"
)

const gopher = `iVBORw0KGgoAAAANSUhEUgAAADIAAAAZCAAAAABWPyKYAAAAAnRSTlMAAHaTzTgAAAAIYWNUTAAAAAgAAAAAuT2L0QAAABpmY1RMAAAAAAAAADIAAAAZAAAAAAAAAAAACgBkAABynX+UAAAA
tklEQVR42o3U16GlQBDEUDKr0BSaQuvd57uAa/SFO/iZY5fPji61tUsASLZNokCZMyFrIYBCHpHgg7ghAW7AfKVyUji/1fF78zZhZKkGgheTmSTiDVm7tIiAc0sGcM7E2d09SZNGKgw3yCZHEr8D9PIObNFEZtyh
qjZJXBWBn1VOBMf7wCL1f70yyZFs4ssSZBteE1jiLZPkPFpek/dHWL+wRvAWaBQuNbgalHRoi7O5zD8J1JZ/S/RPswcO/tkAAAAaZmNUTAAAAAEAAAAyAAAAFwAAAAAAAAAAAAoAZAEAHYa8dAAAALdmZEFUAAAA
AjjLpZNRDsMwCEN3Mx/NR/PRmES7mUFUPvY+2irwhJOorz9B0tdIevFQlIgfKKI5U8nn9abym3hQdOY4JrsjQpVIrsUmQVGwYJoDRhB0TzFckWQHEQGA31k2aqUr4SYrFSsJuG3FwSyBKmcjNud4ZtDNp12DagCw
4v6p1FACNKHv1sG8DU0gxXRuJeFQeFGjgURKCbVCUKk44K60q8FupGC2bLTR/rF9wpw0kIQqvAHYqDphD1TQzAAAABpmY1RMAAAAAwAAAC0AAAAXAAAAAwAAAAAACgBkAQA5L6RYAAAAomZkQVQAAAAEOMuNk9EN
xCAMQ28zj+bRPFruACEbXVT6PuOHSaXyuYBBN+aPnXgoMv2URZxQ7ghm0kMAd9mc7RxmVSmpzRzvfqgSy8HaZ8kEI03ZkbTsKgD0DSk7sl1ObSdhs7q9e3v6CpMQ08+vnADakArytpSNl6FMVkMt5GHnH6V6PoAs
vwKI2OV3m0N+rw/57S6wfW1vniDvbvrgH6p0v5dlPCtH5G3wAAAAGmZjVEwAAAAFAAAAKQAAABkAAAAEAAAAAAAKAGQAAODe53cAAACyZmRBVAAAAAY4y7WTQQ4DMQgD92c8zU/z01IKERZss3vqKIfIjIAccm0s
ucRIlSDw1EyAOxRpkohLQRK8mTwhVeKJWgDkSlp9FbutEVTaREE3bZlhqdDEqOF7DxPmtKYlRg1xcXOH0xzA4kFDRE6cezp7ABMQJBfbM8s0AmSpXmsMU1WZSpsJ3oGfYZ6R+AczVfARyLT3lqXiUbT+PV5ny01K
SH78WJCwDrRiV12cCSr5AIzDOVMiCnshAAAAGmZjVEwAAAAHAAAAKgAAABkAAAAEAAAAAAAKAGQAAK4esjcAAACuZmRBVAAAAAg4y43R0Q3DMAwD0Wx2o3E0jpZGFVI2ClrnPuUHAba3jm5Lmc6ZjgxEnFPBlVK0
D1ALkH2nyL8Cvp38J2Uv2N6TU2a0fE+knAQmPtQgExvpAip7brWkfdBOQGiaNDX1XQaPex1Nh5yLGvJYneROLuYcQKhULoVKdWkunxp8i5Z4kU4qvEw01ZrS9OnSojyTjbWAJYNXKxOtNZsw38s12QpNSGMOl8kL
6SZHxo1/jQoAAAAaZmNUTAAAAAkAAAAwAAAAFwAAAAEAAAAAAAoAZAEAqAUmYwAAALZmZEFUAAAACjjLpZPBDcQwCASvsy1tSqO0nAUn4ayw/Lj5WAEGEpJ8/kACJAsmx3oR6AURwKhUNqqCiDpNsAFxYBCq+/PD
hSfHeXnVuvOKqYXK7GLXdybUQnWQUM8oYc/sQlKHC432W3Lmh2jBjXwlzx4H7VtK3ksF6CvJ9uqQDeIkKCZaAC3mCUzmorTs3uVQaZeAlJYgAlkOCEciUlgGiitKkIYd3T7wu8H8P3Btv9PLwrHyLy61PabHMU88
AAAAGmZjVEwAAAALAAAALQAAABcAAAACAAAAAAAKAGQBAPgYOL4AAACoZmRBVAAAAAw4y4XQ0QnEMAwEUXe2pU1pW5pzWMQSpxjPV1ieIWi8aTXuhQUM8ULRWCvSp8YG4jsYHLm/Q9KfBZ+aBqh2rbtKdxQcqvlZ
U+hYLU+UPLElTxbfeiKJLz31q+qs6dxC9/2tTpSTdCy8b2VL+4LusZx3RePeXQM+JOnwJ7gN0zYagtREbiHB0kbOrv+hq2bhKPUd503u+O5R4vTxgpoEVPsABxlD2YYLabcAAAAaZmNUTAAAAA0AAAAnAAAAGQAA
AAYAAAAAAAoAZAAALHYLBgAAAK1mZEFUAAAADjjLvczBbQMxEEPRdMbSfmksTQFhYrCjwLBP4WFX4jzNT6NXenvXSYBtaUA+JJVljkubdAxO+XCB78IslHyay6SajaJo00eFwnwAX7BHulPZZglruR5nEJdncLZr
MqhTy+VWV7eh4XBDXm5q20BuC1rjDO4/aoflZlo3bdY9Hf4TMIz7kLhAvnT6Z0dcIF+oys+sUiQzbbRUID7WTpqwCyLdDdP8Ap+yPFtXsQ7vAAAAAElFTkSuQmCC`

// gopherPNG creates an io.Reader by decoding the base64 encoded image data string in the gopher constant.
func gopherPNG() io.Reader { return base64.NewDecoder(base64.StdEncoding, strings.NewReader(gopher)) }

func ExampleDecodeAll() {
	am, err := apng.DecodeAll(gopherPNG())
	if err != nil {
		panic(err)
	}

	levels := []rune(" ░▒▓█")
	loops := int64(am.LoopCount)
	if loops == 0 {
		loops = -1
	}
	rect := am.Frames[0].Image.Bounds()
	buf := make([][]rune, rect.Dy())
	for y := range buf {
		buf[y] = make([]rune, rect.Dx())
		for x := range buf[y] {
			buf[y][x] = ' '
		}
	}
	n := time.Now()

	print("\x1b[2J") // clear screen
	for loops != 0 {
		for i, fr := range am.Frames {
			for y := 0; y < fr.Image.Bounds().Dy(); y++ {
				for x := 0; x < fr.Image.Bounds().Dx(); x++ {
					c := fr.Image.At(x, y)
					if _, _, _, a := c.RGBA(); a > 0 {
						buf[y+fr.YOffset][x+fr.XOffset] = levels[(color.GrayModel.Convert(c).(color.Gray).Y / 52)]
					} else if fr.BlendOp == apng.BLEND_OP_SOURCE {
						buf[y+fr.YOffset][x+fr.XOffset] = ' '
					}
				}
			}
			print("\x1b[H") // move top left
			for y := range buf {
				println(string(buf[y]))
			}
			time.Sleep(time.Duration(float64(time.Second)*float64(fr.DelayNumerator)/float64(fr.DelayDenominator)) - time.Now().Sub(n))
			n = time.Now()
			if fr.DisposeOp == apng.DISPOSE_OP_BACKGROUND || i == len(am.Frames)-1 {
				for y := 0; y < fr.Image.Bounds().Dy(); y++ {
					for x := 0; x < fr.Image.Bounds().Dx(); x++ {
						buf[y+fr.YOffset][x+fr.XOffset] = ' '
					}
				}
			}
		}
		if loops > 0 {
			loops--
		}
	}
}

func ExampleEncode() {
	w, h := 200, 200
	am := apng.APNG{Frames: make([]apng.Frame, 20)}
	var circles [3]image.Point

	for i := range am.Frames {
		im := image.NewRGBA(image.Rect(0, 0, w, h))
		am.Frames[i].Image = im

		theta := float64(i) * 2 * math.Pi / float64(len(am.Frames))
		for c := range circles {
			theta0 := float64(c) * 2 * math.Pi / 3
			circles[c].X = int(float64(w) * (0.5 - 0.15*math.Sin(theta0) - 0.1*math.Sin(theta0+theta)))
			circles[c].Y = int(float64(h) * (0.5 - 0.15*math.Cos(theta0) - 0.1*math.Cos(theta0+theta)))
		}

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var rgb [3]uint8
				for c := range circles {
					dx, dy := x-circles[c].X, y-circles[c].Y
					if dx*dx+dy*dy < int(float64(w*h)/20) {
						rgb[c] = 0xff
					}
				}
				im.Set(x, y, color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 0xff})
			}
		}
	}

	f, err := os.Create("rgb.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := apng.Encode(f, am); err != nil {
		panic(err)
	}
}
