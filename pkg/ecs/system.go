package ecs

import "time"

var entityID = 0

type Entity struct {
	id         int
	components []Component
}

type Component interface{}

type System interface {
	AddComponentsFromEntity(entity *Entity) bool
	Update(delta time.Duration)
}

func AddComponent[C Component](e *Entity, c C) {
	e.components = append(e.components, c)
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

type BaseComponent struct{}

type BaseSystem[C Component] struct {
	Components []C
}

func (s *BaseSystem[C]) AddComponentsFromEntity(e *Entity) bool {
	c, ok := GetComponent[C](e)
	if !ok {
		return false
	}

	s.Components = append(s.Components, c)
	return true
}

type Scene struct {
	systems []System
}

func (s *Scene) AddSystem(system System) {
	s.systems = append(s.systems, system)
}

func (s *Scene) AddEntity(entity *Entity) {
	for _, system := range s.systems {
		_ = system.AddComponentsFromEntity(entity)
	}
}

func (s *Scene) Update(delta time.Duration) {
	for _, system := range s.systems {
		system.Update(delta)
	}
}
