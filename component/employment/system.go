package employment

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/ecs"
)

type Component struct {
	ecs.BaseComponent

	CurrentState State
}

type System struct {
	ecs.BaseSystem[*Component]
}

func (s System) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ecs.EntityMsg:
		for _, c := range s.Components {
			if c.Entity.ID() != msg.EntityID {
				continue
			}
			cmds = append(cmds, s.UpdateComponent(c, msg.Msg))
			break
		}
	default:
		for _, c := range s.Components {
			cmd := s.UpdateComponent(c, msg)
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (System) UpdateComponent(c *Component, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case changeStateMsg:
		c.CurrentState = msg
		return c.CurrentState.OnEnter(c)

	case ecs.TickMsg:
		if c.CurrentState == nil {
			return changeState(c.Entity, IdleState{})
		}

		return c.CurrentState.Update(c, msg)

	default:
		return c.CurrentState.Update(c, msg)
	}
}
