package pda

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdbann/forestry/component/physics"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/ecs"
)

type Component struct {
	ecs.BaseComponent

	Client     *client.Client
	Registered bool
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

func (s System) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ecs.EntityMsg:
		for _, c := range s.Components {
			if c.Entity.ID() != msg.EntityID {
				continue
			}
			return s.UpdateComponent(c, msg.Msg)
		}
	case physics.EncounterMsg:
		for _, c := range s.Components {
			if c.Entity.ID() != msg.Source.ID() {
				continue
			}
			return s.UpdateComponent(c, msg)
		}
	default:
	}

	return nil
}

func (System) UpdateComponent(c *Component, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case attemptRegistrationMsg:
		if c.Registered {
			return nil
		}
		return performRegistration(c)
	case RegisterSuccessMsg:
		c.Registered = true
	case physics.EncounterMsg:
		renderComponent, ok := ecs.GetComponent[*render.Component](msg.Encountered)
		if !ok {
			return nil
		}

		return reportDiscovery(c, renderComponent.Position)
	default:
	}
	return nil
}
