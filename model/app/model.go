package app

import (
	"math/rand"

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

type Params struct {
	Rng *rand.Rand
}

func New(params Params) Model {
	return Model{
		Help: help.New(),
		Keys: DefaultKeys,
		World: world.New(world.Params{
			Rng:  params.Rng,
			Size: geo.Size{Width: 64, Height: 24},
		}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.World.Init(), world.AddPerson)
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

	var cmd tea.Cmd
	m.World, cmd = m.World.Update(msg)
	return m, cmd
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
