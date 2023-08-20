package render

import (
	"time"

	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

type Component struct {
	ecs.BaseComponent

	Position geo.Point
	Rune     rune
}

func (c Component) View() string {
	return string(c.Rune)
}

type System struct {
	ecs.BaseSystem[*Component]
}

func (s System) Update(_ time.Duration) {}
