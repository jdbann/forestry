// Package agent provides a component which allows entities to manage their movement based on set behaviours.
//
// Behaviours are not yet implemented so agents will randomly wander around within the provided boundaries.
package agent

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	Rng       *rand.Rand
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

		var (
			step        geo.Vector
			randomSteps = []geo.Vector{
				{X: 0, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: -1},
				{X: -1, Y: 0},
			}
		)

		s.Rng.Shuffle(len(randomSteps), func(i, j int) {
			randomSteps[i], randomSteps[j] = randomSteps[j], randomSteps[i]
		})

		for _, randomStep := range randomSteps {
			if renderComponent.Position.Add(randomStep).WithinSize(s.WorldSize) {
				step = randomStep
				break
			}
		}

		renderComponent.Position = renderComponent.Position.Add(step)
	}

	return nil
}
