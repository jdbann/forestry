package agent

import (
	"math/rand"
	"time"

	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

type Component struct {
	ecs.BaseComponent

	StepFrequency time.Duration
	SinceLastStep time.Duration
}

type System struct {
	ecs.BaseSystem[*Component]

	Rng       *rand.Rand
	WorldSize geo.Size
}

func (s System) Update(delta time.Duration) {
	for _, component := range s.Components {
		component.SinceLastStep += delta
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
}
