// Package geo provides types for simple geometric concepts as well as helpers
// for manipulating them.
package geo

import "math/rand"

type Point struct{ X, Y int }

func (p Point) Add(vec Vector) Point {
	return Point{X: p.X + vec.X, Y: p.Y + vec.Y}
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Within(bounds Rect) bool {
	return p.X >= bounds.Origin.X &&
		p.X < bounds.Size.Width &&
		p.Y >= bounds.Origin.Y &&
		p.Y < bounds.Size.Height
}

func (p Point) WithinSize(bounds Size) bool {
	return p.Within(Rect{Size: bounds})
}

type Rect struct {
	Origin Point
	Size   Size
}

func (r Rect) PointWithin(rng *rand.Rand) Point {
	return Point{
		X: r.Origin.X + rng.Intn(r.Size.Width-r.Origin.X-1),
		Y: r.Origin.Y + rng.Intn(r.Size.Height-r.Origin.Y-1),
	}
}

type Size struct{ Width, Height int }

func (s Size) PointWithin(rng *rand.Rand) Point {
	return Rect{Size: s}.PointWithin(rng)
}

type Vector struct{ X, Y int }
