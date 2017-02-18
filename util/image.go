package util

import "image"

// ImageColor - returns an image with a single pixel
func ImageColor(r, g, b, a uint8) image.Image {
	return &image.RGBA{
		Pix:    []uint8{r, g, b, a},
		Stride: 4,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{1, 1}},
	}
}
