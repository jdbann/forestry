package ecs

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var entityID = 0

type Entity struct {
	id         int
	components []Component
}

func NewEntity(components ...Component) *Entity {
	e := &Entity{}
	for _, component := range components {
		AddComponent(e, component)
	}
	return e
}

func (e Entity) ID() int {
	return e.id
}

type Component interface {
	Init() tea.Cmd
	SetEntity(entity *Entity)
}

type System interface {
	AddComponentsFromEntity(entity *Entity) tea.Cmd
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
}

func AddComponent[C Component](e *Entity, c C) {
	e.components = append(e.components, c)
	c.SetEntity(e)
}

func GetComponent[C Component](e *Entity) (C, bool) {
	for _, component := range e.components {
		if component, ok := component.(C); ok {
			return component, true
		}
	}

	var component C
	return component, false
}

type BaseComponent struct {
	Entity *Entity
}

func (c *BaseComponent) SetEntity(entity *Entity) {
	c.Entity = entity
}

type BaseSystem[C Component] struct {
	Components []C
}

func (s *BaseSystem[C]) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range s.Components {
		cmds = append(cmds, c.Init())
	}
	return tea.Batch(cmds...)
}

func (s *BaseSystem[C]) AddComponentsFromEntity(e *Entity) tea.Cmd {
	c, ok := GetComponent[C](e)
	if !ok {
		return nil
	}

	s.Components = append(s.Components, c)
	return c.Init()
}

type Scene struct {
	systems []System
}

func (s *Scene) AddSystem(system System) {
	s.systems = append(s.systems, system)
}

func (s *Scene) AddEntity(entity *Entity) tea.Cmd {
	var cmds []tea.Cmd
	for _, system := range s.systems {
		cmds = append(cmds, system.AddComponentsFromEntity(entity))
	}
	return tea.Batch(cmds...)
}

func (s *Scene) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, system := range s.systems {
		cmds = append(cmds, system.Init())
	}
	return tea.Batch(cmds...)
}

func (s *Scene) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	for _, system := range s.systems {
		cmd = system.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

type TickMsg time.Duration
