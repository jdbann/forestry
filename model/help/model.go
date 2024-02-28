package help

import (
	helpbubble "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/pkg/geo"
)

var helpStyle = lipgloss.NewStyle()

type Model struct {
	bubble helpbubble.Model
	keys   KeyMap
	size   geo.Size
}

type Params struct {
	Keys KeyMap
}

func New(params Params) Model {
	return Model{
		keys: params.Keys,
	}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = geo.Size(msg)
		m.bubble.Width = msg.Width
	default:
	}
	return m, nil
}

func (m Model) View() string {
	return helpStyle.Width(m.size.Width).Height(m.size.Height).Render(m.bubble.View(m.keys))
}

type KeyMap struct {
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

var DefaultKeys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
