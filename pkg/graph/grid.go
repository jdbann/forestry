package graph

import (
	"math"

	"github.com/jdbann/forestry/pkg/geo"
)

// GridGraph holds a graph that represents a two dimensional grid, where individual nodes can be identified based on their X and Y coordinates.
type GridGraph struct {
	graph         *graph[geo.Point]
	width, height int
}

// NewGridGraph builds a new GridGraph of the specified size.
func NewGridGraph(width, height int) *GridGraph {
	g := new[geo.Point]()

	nodeID := func(x, y int) int { return x + (y * width) }

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			g.addNode(nodeID(x, y), geo.Point{X: x, Y: y})

			if x > 0 {
				g.addEdge(nodeID(x, y), nodeID(x-1, y), 1)
				g.addEdge(nodeID(x-1, y), nodeID(x, y), 1)
			}

			if y > 0 {
				g.addEdge(nodeID(x, y), nodeID(x, y-1), 1)
				g.addEdge(nodeID(x, y-1), nodeID(x, y), 1)
			}
		}
	}

	return &GridGraph{
		graph:  g,
		width:  width,
		height: height,
	}
}

// FindPath attempts to find a path from the starting location to the ending location. It returns the list of points that make up the requested path, and whether or not a path could be successfully found.
func (gg *GridGraph) FindPath(from, to geo.Point) ([]geo.Point, bool) {
	nodeID := func(x, y int) int { return x + (y * gg.width) }

	nodes, ok := gg.graph.findPath(nodeID(from.X, from.Y), nodeID(to.X, to.Y), gridHeuristic)
	if !ok {
		return nil, false
	}

	path := []geo.Point{}
	for _, node := range nodes {
		path = append(path, gg.graph.nodes[node].val)
	}

	return path, true
}

// FindNeighbours returns the coordinates for the nodes that are neighbours of the target node.
func (gg *GridGraph) FindNeighbours(target geo.Point) []geo.Point {
	nodeID := func(x, y int) int { return x + (y * gg.width) }

	neighbourIDs := gg.graph.neighbours(nodeID(target.X, target.Y))
	neighbours := make([]geo.Point, len(neighbourIDs))

	for i, id := range neighbourIDs {
		neighbours[i] = gg.graph.nodes[id].val
	}

	return neighbours
}

func (gg *GridGraph) Size() geo.Size {
	return geo.Size{Width: gg.width, Height: gg.height}
}

func gridHeuristic(p1, p2 geo.Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2))
}
