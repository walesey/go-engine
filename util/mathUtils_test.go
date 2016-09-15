package util

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
)

func expectApproxEqual2D(t *testing.T, v1, v2 mgl32.Vec2) {
	assert.True(t, v1.ApproxEqualThreshold(v2, 0.001), fmt.Sprintf("Expected %v to approx equal %v", v1, v2))
}

func TestRound(t *testing.T) {
	var n float32 = 1.999
	i := Round(n, .5, 0)
	assert.EqualValues(t, 2.0, i, "LengthSquared")
}

func TestRoundHalfUp(t *testing.T) {
	var n float32 = 1.501
	i := RoundHalfUp(n)
	assert.EqualValues(t, 2, i, "LengthSquared")
}

func TestPointToLineDist(t *testing.T) {
	assert.EqualValues(t, 1, PointToLineDist(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 0, 1}, mgl32.Vec3{0, 1, 0}))
}

func TestTwoSegmentIntersect_intersection(t *testing.T) {
	intersect, err := TwoSegmentIntersect(mgl32.Vec2{0, 0}, mgl32.Vec2{2, 2}, mgl32.Vec2{-1, 1}, mgl32.Vec2{2, 1})
	assert.Nil(t, err, "TwoSegmentIntersect error should be nil")
	assert.EqualValues(t, mgl32.Vec2{1, 1}, intersect, "TwoSegmentIntersect intersect")
}

func TestTwoSegmentIntersect_noIntersection(t *testing.T) {
	_, err := TwoSegmentIntersect(mgl32.Vec2{0, 0}, mgl32.Vec2{-2, -2}, mgl32.Vec2{-1, 1}, mgl32.Vec2{2, 1})
	assert.NotNil(t, err, "TwoSegmentIntersect error should not be nil")
}

func TestTwoSegmentIntersect_largeNumbers(t *testing.T) {
	intersect, err := TwoSegmentIntersect(mgl32.Vec2{455, 494}, mgl32.Vec2{974, 76}, mgl32.Vec2{835, 174}, mgl32.Vec2{961, 211})
	assert.Nil(t, err, "TwoSegmentIntersect error should be nil")
	expectApproxEqual2D(t, mgl32.Vec2{847.693, 177.727}, intersect)
}

func TestSegmentCircleIntersect_intersectionEntry(t *testing.T) {
	intersect, err := SegmentCircleIntersect(0.5, mgl32.Vec2{2, 1}, mgl32.Vec2{1, 3}, mgl32.Vec2{2, 1})
	assert.Nil(t, err, "SegmentCircleIntersect error should be nil")
	expectApproxEqual2D(t, mgl32.Vec2{1.776, 1.447}, intersect)
}

func TestSegmentCircleIntersect_intersectionExit(t *testing.T) {
	intersect, err := SegmentCircleIntersect(0.5, mgl32.Vec2{2, 1}, mgl32.Vec2{2, 1}, mgl32.Vec2{1, 3})
	assert.Nil(t, err, "SegmentCircleIntersect error should be nil")
	expectApproxEqual2D(t, mgl32.Vec2{1.776, 1.447}, intersect)
}

func TestSegmentCircleIntersect_2intersections(t *testing.T) {
	intersect, err := SegmentCircleIntersect(0.5, mgl32.Vec2{2, 1}, mgl32.Vec2{1, 3}, mgl32.Vec2{4, -2})
	assert.Nil(t, err, "SegmentCircleIntersect error should be nil")
	expectApproxEqual2D(t, mgl32.Vec2{1.905, 1.491}, intersect)
}

func TestSegmentCircleIntersect_noIntersection(t *testing.T) {
	_, err := SegmentCircleIntersect(0.5, mgl32.Vec2{2, 1}, mgl32.Vec2{1, 1}, mgl32.Vec2{5, -2})
	assert.NotNil(t, err, "SegmentCircleIntersect error should not be nil")
}
