package brain

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/ecs"
)

type Component struct {
	ecs.BaseComponent

	AvailableStates []State
	CurrentState    State
}

func (c *Component) Init() tea.Cmd {
	return c.CurrentState.OnEnter(c)
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(brainMsg); ok && msg.id != c.Entity.ID() {
		return nil
	}

	cmd := c.CurrentState.Update(c, msg)
	if cmd == nil {
		return nil
	}

	return wrapCmd(c.Entity.ID(), cmd)
}

func (c *Component) EnterState(target State) tea.Cmd {
	for _, s := range c.AvailableStates {
		if s == target {
			c.CurrentState = s
			return s.OnEnter(c)
		}
	}

	return nil
}

type System struct {
	ecs.BaseSystem[*Component]
}

func (s System) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for _, s := range s.Components {
		cmd = s.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

type State interface {
	Update(*Component, tea.Msg) tea.Cmd
	OnEnter(*Component) tea.Cmd
}

type BaseState struct{}

func (s BaseState) Update(_ *Component, _ tea.Msg) tea.Cmd {
	return nil
}

func (s BaseState) OnEnter(_ *Component) tea.Cmd {
	return nil
}

type brainMsg struct {
	id  int
	msg tea.Msg
}

func wrapCmd(id int, cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		msg := cmd()
		return brainMsg{
			id:  id,
			msg: msg,
		}
	}
}
