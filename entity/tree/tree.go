package tree

import (
	"github.com/jdbann/forestry/component/physics"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

func New(at geo.Point) *ecs.Entity {
	return ecs.NewEntity(
		&render.Component{
			Position: at,
			Rune:     '∆',
		},
		&physics.Component{},
	)
}
