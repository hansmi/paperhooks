package main

import (
	"bytes"
	"image/png"
	"testing"
)

func TestMakeRandomImage(t *testing.T) {
	b, err := makeRandomImage(1, 1)
	if err != nil {
		t.Fatalf("makeRandomImage() failed: %v", err)
	}

	if _, err := png.Decode(bytes.NewReader(b)); err != nil {
		t.Errorf("Decoding image failed: %v", err)
	}
}
