// Copyright 2019 kts of kettek / Ketchetwahmeegwun Tecumseh Southall. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apng

import (
	"os"
	"testing"
)

func TestReadAPNGWithDefaultFrame(t *testing.T) {
	a, err := ReadAPNG("tests/WithDefaultFrame.png")
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
	a, err := ReadAPNG("tests/WithoutDefaultFrame.png")
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
	a, err := ReadAPNG("tests/MultipleIDATs.png")
	if err != nil {
		t.Error(err)
		return
	}

	if len(a.Frames) != 2 {
		t.Error("Expected 2 frames.")
		return
	}
}

func ReadAPNG(path string) (APNG, error) {
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
