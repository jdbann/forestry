package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/model/world"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/geo"
)

type Model struct {
	Size  geo.Size
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

	case tea.WindowSizeMsg:
		m.Size = geo.Size{Width: msg.Width, Height: msg.Height}
	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.Place(
		m.Size.Width, m.Size.Height,
		lipgloss.Center, lipgloss.Center,
		m.World.View(),
		lipgloss.WithWhitespaceBackground(color.Gray1),
	)
}
