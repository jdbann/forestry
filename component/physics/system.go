package physics

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/component/agent"
	"github.com/jdbann/forestry/component/graph"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
)

type Component struct {
	ecs.BaseComponent
}

type System struct {
	ecs.BaseSystem[*Component]
}

func (s System) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ecs.EntityMsg:
		for _, c := range s.Components {
			if c.Entity.ID() != msg.EntityID {
				continue
			}

			return s.UpdateComponent(c, msg.Msg)
		}
	default:
	}

	return nil
}

func (s System) UpdateComponent(c *Component, msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.(type) {
	case agent.EntityMovedMsg:
		renderComponent, ok := ecs.GetComponent[*render.Component](c.Entity)
		if !ok {
			return nil
		}

		graphComponent, ok := ecs.GetComponent[*graph.Component](c.Entity)
		if !ok {
			return nil
		}

		for _, neighbour := range graphComponent.Graph.FindNeighbours(renderComponent.Position) {
			for _, otherComponent := range s.Components {
				if c.Entity.ID() == otherComponent.Entity.ID() {
					continue
				}

				otherRenderComponent, ok := ecs.GetComponent[*render.Component](otherComponent.Entity)
				if !ok {
					break
				}

				if neighbour.Equals(otherRenderComponent.Position) {
					cmds = append(cmds, ReportEncountered(c.Entity, otherComponent.Entity))
				}
			}
		}

		return tea.Batch(cmds...)
	default:
	}
	return nil
}
