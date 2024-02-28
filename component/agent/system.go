// Package agent provides a component which allows entities to manage their movement based on set behaviours.
//
// Behaviours are not yet implemented so agents will randomly wander around within the provided boundaries.
package agent

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/component/graph"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

// Component keeps track of an entity's speed and determines if enough time has elapsed for the next step to be taken.
type Component struct {
	ecs.BaseComponent

	StepFrequency time.Duration
	SinceLastStep time.Duration
	Path          []geo.Point
}

// System manages agent components and applies each component's behaviour.
type System struct {
	ecs.BaseSystem[*Component]

	WorldSize geo.Size
}

// Update handles tick messages, moving agents if enough time has elapsed.
func (s System) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case agentMsg:
		for _, c := range s.Components {
			if c.Entity.ID() != msg.id {
				continue
			}

			s.UpdateComponent(c, msg.msg)
		}
	default:
		for _, c := range s.Components {
			s.UpdateComponent(c, msg)
		}
	}

	return nil
}

func (System) UpdateComponent(c *Component, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case setDestinationMsg:
		renderComponent, ok := ecs.GetComponent[*render.Component](c.Entity)
		if !ok {
			return nil
		}

		graphComponent, ok := ecs.GetComponent[*graph.Component](c.Entity)
		if !ok {
			return nil
		}

		path, ok := graphComponent.Graph.FindPath(renderComponent.Position, geo.Point(msg))
		if !ok {
			return nil
		}

		c.Path = path
		return nil
	case ecs.TickMsg:
		c.SinceLastStep += time.Duration(msg)
		if c.SinceLastStep < c.StepFrequency {
			return nil
		}

		c.SinceLastStep -= c.StepFrequency

		if len(c.Path) == 0 {
			return nil
		}

		renderComponent, ok := ecs.GetComponent[*render.Component](c.Entity)
		if !ok {
			return nil
		}

		renderComponent.Position, c.Path = c.Path[0], c.Path[1:]
	}

	return nil
}

type agentMsg struct {
	id  int
	msg tea.Msg
}

func wrapCmd(id int, cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	return func() tea.Msg {
		return wrapMsg(id, cmd())
	}
}

func wrapMsg(id int, msg tea.Msg) tea.Msg {
	return agentMsg{
		id:  id,
		msg: msg,
	}
}
