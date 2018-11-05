# APNG golang library
This `apng` package provides methods for decoding APNG files.

If a regular PNG file is read, the first Frame as returned by `Decode(*File)` will be the PNG data.

## Methods
### DecodeAll(*File) (APNG, error)
This method returns an APNG type containing the frames and associated data within the passed file.

### DecodeFirst(*File) (image.Image, error)
This method returns the Image of the default frame of an APNG file.

## Types
### APNG
The APNG type contains the frames of a decoded `.apng` file, along with any important properties. It may also be created and used for Encoding.

| Signature                 | Description                   |
|---------------------------|-------------------------------|
| Frames []Frame            | The stored frames of the APNG.|
| LoopCount uint            | The number of times an animation should be restarted during display. A value of 0 means to loop forever.   |

### Frame
The Frame type contains an individual frame of an APNG. The following table provides the important properties and methods.

| Signature                 | Description      |
|---------------------------|------------------|
| Img image.Image           | Frame image data. |
| IsDefault bool            | Indicates if this frame is a default image that should not be included as part of the animation frames. May only be true for the first Frame. |
| XOffset int               | Returns the x offset of the frame. |
| YOffset int               | Returns the y offset of the frame. |
| DelayNumerator int        | Returns the delay numerator.       |
| DelayDenominator int      | Returns the delay denominator.     |
| DisposalOp byte           | Returns the frame disposal operation. May be `apng.DISPOSE_OP_NONE`, `apng.DISPOSE_OP_BACKGROUND`, or `apng.DISPOSE_OP_PREVIOUS`. See the [APNG Specification](https://wiki.mozilla.org/APNG_Specification#.60fcTL.60:_The_Frame_Control_Chunk) for more information. |
| BlendOp byte              | Returns the frame blending operation. May be `apng.BLEND_OP_SOURCE` or `apng.BLEND_OP_OVER`. See the [APNG Specification](https://wiki.mozilla.org/APNG_Specification#.60fcTL.60:_The_Frame_Control_Chunk) for more information. |

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
  a, err := apng.DecodeAll(f)
  if err != nil {
    panic(err)
  }

  log.Printf("Found %d frames\n", len(a.Frames))
  for i, frame := range a.Frames {
    b := frame.Img.Bounds()
    log.Printf("Frame %d: %dx%d\n", i, b.Max.X, b.Max.Y)
  }
}

```
