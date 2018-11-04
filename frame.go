// Copyright 2018 kts of kettek / Ketchetwahmeegwun Tecumseh Southall. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apng

import (
	"image"
)

// dispose_op values, as per the APNG spec.
const (
	DISPOSE_OP_NONE       = 0
	DISPOSE_OP_BACKGROUND = 1
	DISPOSE_OP_PREVIOUS   = 2
)

// blend_op values, as per the APNG spec.
const (
	BLEND_OP_SOURCE = 0
	BLEND_OP_OVER   = 1
)

type Frame struct {
	Img                image.Image
	width, height      int
	x_offset, y_offset int
	delay_num          uint16
	delay_den          uint16
	dispose_op         byte
	blend_op           byte
	is_default         bool
}

// IsDefault indicates if the Frame is a default image that
// should not be used in the animation. IsDefault() may only
// return true on the first frame.
func (f *Frame) IsDefault() bool {
	return f.is_default
}
func (f *Frame) GetWidth() int {
	return f.width
}
func (f *Frame) GetHeight() int {
	return f.height
}
func (f *Frame) GetXOffset() int {
	return f.x_offset
}
func (f *Frame) GetYOffset() int {
	return f.y_offset
}
func (f *Frame) GetDelayNumerator() int {
	return int(f.delay_num)
}
func (f *Frame) GetDelayDenominator() int {
	return int(f.delay_den)
}
func (f *Frame) GetDisposal() byte {
	return f.dispose_op
}
func (f *Frame) GetBlend() byte {
	return f.blend_op
}
