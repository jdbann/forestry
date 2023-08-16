package world

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/geo"
)

var mapStyle = lipgloss.NewStyle().Background(color.Grass3).Foreground(color.Grass11)

type Model struct {
	Size geo.Size
}

func New(size geo.Size) Model {
	return Model{
		Size: size,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return mapStyle.Width(m.Size.Width).Height(m.Size.Height).Render()
}
