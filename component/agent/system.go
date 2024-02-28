// Package agent provides a component which allows entities to manage their movement based on set behaviours.
//
// Behaviours are not yet implemented so agents will randomly wander around within the provided boundaries.
package agent

import (
	"math/rand"
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
}

// System manages agent components and applies each component's behaviour.
type System struct {
	ecs.BaseSystem[*Component]

	WorldSize geo.Size
}

// Update handles tick messages, moving agents if enough time has elapsed.
func (s System) Update(msg tea.Msg) tea.Cmd {
	delta, ok := msg.(ecs.TickMsg)
	if !ok {
		return nil
	}

	for _, component := range s.Components {
		component.SinceLastStep += time.Duration(delta)
		if component.SinceLastStep < component.StepFrequency {
			continue
		}

		component.SinceLastStep -= component.StepFrequency

		renderComponent, ok := ecs.GetComponent[*render.Component](component.Entity)
		if !ok {
			continue
		}

		graphComponent, ok := ecs.GetComponent[*graph.Component](component.Entity)
		if !ok {
			continue
		}

		neighbours := graphComponent.Graph.FindNeighbours(renderComponent.Position)

		if len(neighbours) == 0 {
			continue
		}

		rand.Shuffle(len(neighbours), func(i, j int) {
			neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
		})

		renderComponent.Position = neighbours[0]
	}

	return nil
}
