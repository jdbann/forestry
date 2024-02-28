package employment

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/component/agent"
	"github.com/jdbann/forestry/component/graph"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

type State interface {
	OnEnter(c *Component) tea.Cmd
	Update(c *Component, msg tea.Msg) tea.Cmd
}

type IdleState struct{}

func (IdleState) OnEnter(c *Component) tea.Cmd {
	renderComponent, ok := ecs.GetComponent[*render.Component](c.Entity)
	if !ok {
		return nil
	}

	renderComponent.Rune = 'I'

	graphComponent, ok := ecs.GetComponent[*graph.Component](c.Entity)
	if !ok {
		return nil
	}

	destination := graphComponent.Graph.Size().RandomPointWithin()

	return tea.Sequence(
		sleepCmd(time.Millisecond*time.Duration(rand.Intn(2000))),
		tea.Batch(
			agent.SetDestination(c.Entity, destination),
			changeState(c.Entity, WalkingState{destination: destination}),
		),
	)
}

func (IdleState) Update(_ *Component, _ tea.Msg) tea.Cmd {
	return nil
}

type WalkingState struct {
	destination geo.Point
}

func (WalkingState) OnEnter(c *Component) tea.Cmd {
	renderComponent, ok := ecs.GetComponent[*render.Component](c.Entity)
	if !ok {
		return nil
	}

	renderComponent.Rune = 'W'
	return nil
}

func (s WalkingState) Update(c *Component, msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case ecs.TickMsg:
		renderComponent, ok := ecs.GetComponent[*render.Component](c.Entity)
		if !ok {
			return nil
		}

		if renderComponent.Position.Equals(s.destination) {
			return changeState(c.Entity, IdleState{})
		}
	default:
	}
	return nil
}

func sleepCmd(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return nil
	}
}
