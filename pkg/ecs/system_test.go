package ecs_test

import (
	"testing"
	"time"

	"github.com/jdbann/forestry/pkg/ecs"
	"gotest.tools/v3/assert"
)

func TestEntityComponentSystem(t *testing.T) {
	scene := &ecs.Scene{}
	system := &countSystem{}
	scene.AddSystem(system)
	entity := &ecs.Entity{}

	runStep(t, "try to get component when not added to entity", func(t *testing.T) {
		component, ok := ecs.GetComponent[*countComponent](entity)
		assert.Assert(t, component == nil)
		assert.Assert(t, !ok)
	})

	scene.AddEntity(entity)

	runStep(t, "get component after adding to entity", func(t *testing.T) {
		ecs.AddComponent[*countComponent](entity, &countComponent{})
		component, ok := ecs.GetComponent[*countComponent](entity)
		assert.Assert(t, component != nil)
		assert.Assert(t, ok)
	})

	runStep(t, "add component from entity to system", func(t *testing.T) {
		ok := system.AddComponentsFromEntity(entity)
		assert.Assert(t, ok)
	})

	runStep(t, "update system and confirm component was updated", func(t *testing.T) {
		scene.Update(0)
		component, ok := ecs.GetComponent[*countComponent](entity)
		assert.Assert(t, ok)
		assert.Equal(t, component.count, 1)
	})
}

type countComponent struct {
	ecs.BaseComponent
	count int
}

type countSystem struct {
	ecs.BaseSystem[*countComponent]
}

func (s *countSystem) Update(_ time.Duration) {
	for _, component := range s.BaseSystem.Components {
		component.count++
	}
}

func runStep(t *testing.T, name string, fn func(t *testing.T)) {
	if !t.Run(name, fn) {
		t.FailNow()
	}
}
