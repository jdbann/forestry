package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/model/world"
	"github.com/jdbann/forestry/pkg/geo"
)

type Model struct {
	Help  help.Model
	Keys  KeyMap
	Size  geo.Size
	World tea.Model
}

func New() Model {
	return Model{
		Help:  help.New(),
		Keys:  DefaultKeys,
		World: world.New(geo.Size{Width: 64, Height: 24}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Size = geo.Size{Width: msg.Width, Height: msg.Height}
		m.Help.Width = m.Size.Width
	}

	return m, nil
}

func (m Model) View() string {
	helpView := m.Help.View(m.Keys)
	worldView := lipgloss.Place(
		m.Size.Width, m.Size.Height-lipgloss.Height(helpView),
		lipgloss.Center, lipgloss.Center,
		m.World.View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, worldView, helpView)
}
