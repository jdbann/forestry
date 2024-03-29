package agent

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

func SetDestination(e *ecs.Entity, dest geo.Point) tea.Cmd {
	return func() tea.Msg {
		return ecs.EntityMsg{
			EntityID: e.ID(),
			Msg:      setDestinationMsg(dest),
		}
	}
}

type setDestinationMsg geo.Point
