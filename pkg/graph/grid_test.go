package graph_test

import (
	"testing"

	"github.com/jdbann/forestry/pkg/geo"
	"github.com/jdbann/forestry/pkg/graph"
	"gotest.tools/v3/assert"
)

func TestGridGraph(t *testing.T) {
	type testCase struct {
		name          string
		width, height int
		from, to      geo.Point
		expectPath    []geo.Point
		expectOK      bool
	}

	run := func(t *testing.T, tc testCase) {
		g := graph.NewGridGraph(tc.width, tc.height)

		path, ok := g.FindPath(tc.from, tc.to)
		assert.Equal(t, tc.expectOK, ok)
		assert.DeepEqual(t, tc.expectPath, path)
	}

	testCases := []testCase{
		{
			name:   "3x3 grid",
			width:  3,
			height: 3,
			from:   geo.Point{X: 0, Y: 0},
			to:     geo.Point{X: 2, Y: 2},
			expectPath: []geo.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 1, Y: 2},
				{X: 2, Y: 2},
			},
			expectOK: true,
		},
		{
			name:       "3x3 grid with destination outside bounds",
			width:      3,
			height:     3,
			from:       geo.Point{X: 0, Y: 0},
			to:         geo.Point{X: 9, Y: 9},
			expectPath: nil,
			expectOK:   false,
		},
		{
			name:   "10x10 grid",
			width:  10,
			height: 10,
			from:   geo.Point{X: 4, Y: 5},
			to:     geo.Point{X: 5, Y: 0},
			expectPath: []geo.Point{
				{X: 4, Y: 5},
				{X: 4, Y: 4},
				{X: 4, Y: 3},
				{X: 4, Y: 2},
				{X: 4, Y: 1},
				{X: 4, Y: 0},
				{X: 5, Y: 0},
			},
			expectOK: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
