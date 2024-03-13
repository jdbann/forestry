package person

import (
	"time"

	"github.com/jdbann/forestry/component/agent"
	"github.com/jdbann/forestry/component/employment"
	"github.com/jdbann/forestry/component/graph"
	"github.com/jdbann/forestry/component/pda"
	"github.com/jdbann/forestry/component/physics"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

func New(at geo.Point) *ecs.Entity {
	e := ecs.NewEntity(
		&render.Component{
			Position: at,
			Rune:     'P',
		},
		&agent.Component{
			StepFrequency: time.Millisecond * 100,
			SinceLastStep: 0,
		},
		&pda.Component{},
		&graph.Component{},
		&employment.Component{},
		&physics.Component{},
	)
	return e
}
