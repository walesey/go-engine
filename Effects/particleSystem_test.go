package effects

import (
    "testing"
	// "image"
	"image/color"

    "github.com/stretchr/testify/assert"
)

func TestLerpColor(t *testing.T) {
	c := lerpColor( color.NRGBA{123,10,55,50}, color.NRGBA{223,60,20,240}, 0.0 )
    assert.EqualValues(t, color.NRGBA{123,10,55,50}, c, "Lerp Color")

	c = lerpColor( color.NRGBA{123,10,55,50}, color.NRGBA{223,60,20,240}, 1.0 )
    assert.EqualValues(t, color.NRGBA{223,60,20,240}, c, "Lerp Color")

	c = lerpColor( color.NRGBA{120,10,50,50}, color.NRGBA{220,60,20,240}, 0.5 )
    assert.EqualValues(t, color.NRGBA{170,35,35,145}, c, "Lerp Color")
}