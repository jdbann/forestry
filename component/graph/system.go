package graph

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
	"github.com/jdbann/forestry/pkg/graph"
)

type Component struct {
	ecs.BaseComponent

	Graph *graph.GridGraph
}

type System struct {
	ecs.BaseSystem[*Component]

	Graph *graph.GridGraph
}

func NewSystem(size geo.Size) *System {
	return &System{
		Graph: graph.NewGridGraph(size.Width, size.Height),
	}
}

func (s *System) Init() tea.Cmd {
	for _, c := range s.Components {
		c.Graph = s.Graph
	}

	return nil
}

func (s *System) AddComponentsFromEntity(e *ecs.Entity) tea.Cmd {
	c, ok := ecs.GetComponent[*Component](e)
	if !ok {
		return nil
	}

	c.Graph = s.Graph

	s.Components = append(s.Components, c)
	return c.Init()
}

func (System) Update(_ tea.Msg) tea.Cmd {
	return nil
}
