package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math/rand"
)

func makeRandomImage(width, height int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	if _, err := rand.Read(img.Pix); err != nil {
		return nil, fmt.Errorf("generating random image: %w", err)
	}

	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}

	return buf.Bytes(), nil
}
