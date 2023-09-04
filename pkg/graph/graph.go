package graph

import (
	"github.com/jdbann/forestry/pkg/priority"
)

type graph[V any] struct {
	nodes map[int]*node[V]
	edges map[int][]*edge[V]
}

type node[V any] struct {
	id  int
	val V
}

type edge[V any] struct {
	weight float64
	node   *node[V]
}

func new[V any]() *graph[V] {
	return &graph[V]{
		nodes: make(map[int]*node[V]),
		edges: make(map[int][]*edge[V]),
	}
}

func (g *graph[V]) addNode(id int, val V) {
	g.nodes[id] = &node[V]{id: id, val: val}
}

func (g *graph[V]) addEdge(from, to int, weight float64) {
	if _, ok := g.nodes[from]; !ok {
		return
	}

	dest, ok := g.nodes[to]
	if !ok {
		return
	}

	g.edges[from] = append(g.edges[from], &edge[V]{weight: weight, node: dest})
}

func (g *graph[V]) neighbours(id int) []int {
	result := []int{}

	for _, edge := range g.edges[id] {
		result = append(result, edge.node.id)
	}

	return result
}

func (g *graph[V]) edge(from, to int) *edge[V] {
	edges, ok := g.edges[from]
	if !ok {
		return nil
	}

	for _, edge := range edges {
		if edge.node.id == to {
			return edge
		}
	}

	return nil
}

func (g *graph[V]) findPath(start, end int, heuristic func(V, V) float64) ([]int, bool) {
	if g.nodes[start] == nil || g.nodes[end] == nil {
		return nil, false
	}

	if heuristic == nil {
		heuristic = func(v1, v2 V) float64 { return 0 }
	}

	// Track untraversed nodes and their priorities (cost + heuristic)
	frontier := priority.NewQueue[int](len(g.nodes))
	frontier.Push(start, 0)

	type step struct {
		node int
		cost float64
	}

	// Track lowest cost routes to traversed nodes
	cameFrom := make(map[int]step, len(g.nodes))
	cameFrom[start] = step{start, 0}

	var current int

	for frontier.Len() > 0 {
		current = frontier.Pop()

		if current == end {
			break
		}

		// Traverse all edges of current node
		for _, next := range g.neighbours(current) {
			// Cost to reach neighbour is current cost + weight of edge
			newCost := cameFrom[current].cost + g.edge(current, next).weight

			// If node can already be reached by lower cost path, skip this edge
			cf, ok := cameFrom[next]
			if ok && newCost >= cf.cost {
				continue
			}

			// Record the shortest route to reach next node
			cameFrom[next] = step{current, newCost}

			// Add next node to frontier
			frontier.Push(next, newCost+heuristic(g.nodes[next].val, g.nodes[end].val))
		}
	}

	if _, ok := cameFrom[end]; !ok {
		return nil, false
	}

	// Assemble path by traversing cameFrom map back to origin
	path := []int{end}
	for path[0] != start {
		cf := cameFrom[path[0]]
		path = append([]int{cf.node}, path...)
	}

	return path, true
}
