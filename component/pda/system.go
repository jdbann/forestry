package pda

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/ecs"
)

type Component struct {
	ecs.BaseComponent

	Client *client.Client
}

type System struct {
	ecs.BaseSystem[*Component]

	Client *client.Client
}

func (s *System) Init() tea.Cmd {
	for _, c := range s.Components {
		c.Client = s.Client
	}

	return nil
}

func (s *System) AddComponentsFromEntity(e *ecs.Entity) tea.Cmd {
	c, ok := ecs.GetComponent[*Component](e)
	if !ok {
		return nil
	}

	c.Client = s.Client

	s.Components = append(s.Components, c)
	return c.Init()
}

func (System) Update(_ tea.Msg) tea.Cmd {
	return nil
}
