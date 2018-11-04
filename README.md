# APNG golang library
This `apng` package provides methods for decoding APNG files.

If a regular PNG file is read, the first Frame as returned by `Decode(*File)` will be the PNG data.

## Methods
### Decode(*File) []Frame
This method returns a slice of `Frame`s.

### DecodeFirst(*File) Image
This method returns the Image of the first `Frame` of an APNG.

## Types
### Frame
The Frame type contains an individual frame of an APNG. The following table provides the important properties and methods.

| Signature                 | Description      |
|---------------------------|------------------|
| Img image.Image           | Frame image data. |
| IsDefault() bool          | Indicates if this frame is a default image that should not be included as part of the animation frames. May only be true for the first Frame. |
| GetWidth() int            | Returns the width of the frame.    |
| GetHeight() int           | Returns the height of the frame.   |
| GetXOffset() int          | Returns the x offset of the frame. |
| GetYOffset() int          | Returns the y offset of the frame. |
| GetDelayNumerator() int   | Returns the delay numerator.       |
| GetDelayDenominator() int | Returns the delay denominator.     |
| GetDisposal() byte        | Returns the frame disposal operation. May be `apng.DISPOSE_OP_NONE`, `apng.DISPOSE_OP_BACKGROUND`, or `apng.DISPOSE_OP_PREVIOUS`. See the [APNG Specification](https://wiki.mozilla.org/APNG_Specification#.60fcTL.60:_The_Frame_Control_Chunk) for more information. |
| GetBlend() byte           | Returns the frame blending operation. May be `apng.BLEND_OP_SOURCE` or `apng.BLEND_OP_OVER`. See the [APNG Specification](https://wiki.mozilla.org/APNG_Specification#.60fcTL.60:_The_Frame_Control_Chunk) for more information. |

## Example
```go
import (
  "github.com/kettek/apng"
  "os"
  "log"
)

func main() {
  f, err := os.Open("my_file.png")
  if err != nil {
    panic(err)
  }
  defer f.Close()
  frames, err := apng.Decode(f)
  if err != nil {
    panic(err)
  }

  log.Printf("Found %d frames\n", len(frames))
  for i, frame := range frames {
    log.Printf("Frame %d: %dx%d\n", i, frame.GetWidth(), frame.GetHeight())
  }
}

```
