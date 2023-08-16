package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/model/world"
	"github.com/jdbann/forestry/pkg/geo"
)

type Model struct {
	World tea.Model
}

func New() Model {
	return Model{
		World: world.New(geo.Size{Width: 64, Height: 24}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	return m.World.View()
}
