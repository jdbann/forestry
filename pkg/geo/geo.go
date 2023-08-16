// Package geo provides types for simple geometric concepts as well as helpers
// for manipulating them.
package geo

type Point struct{ X, Y int }

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

type Size struct{ Width, Height int }
