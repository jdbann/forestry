package physics

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/ecs"
)

type EncounterMsg struct {
	Source      *ecs.Entity
	Encountered *ecs.Entity
}

func ReportEncountered(source, encountered *ecs.Entity) tea.Cmd {
	return func() tea.Msg {
		return EncounterMsg{
			Source:      source,
			Encountered: encountered,
		}
	}
}
