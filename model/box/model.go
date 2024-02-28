package box

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/pkg/geo"
)

type Model struct {
	Size  geo.Size
	Style lipgloss.Style
	Text  string
}

type Params struct {
	Fill lipgloss.TerminalColor
	Size geo.Size
	Text string
}

func New(params Params) Model {
	return Model{
		Size:  params.Size,
		Style: lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Background(params.Fill),
		Text:  params.Text,
	}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = geo.Size(msg)
	default:
	}
	return m, nil
}

func (m Model) View() string {
	return m.Style.Width(m.Size.Width).Height(m.Size.Height).Render(m.Text)
}
