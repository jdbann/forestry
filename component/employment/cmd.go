package employment

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/ecs"
)

func changeState(e *ecs.Entity, nextState State) tea.Cmd {
	return func() tea.Msg {
		return ecs.EntityMsg{
			EntityID: e.ID(),
			Msg:      changeStateMsg(nextState),
		}
	}
}

type changeStateMsg State
