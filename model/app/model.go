package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdbann/forestry/model/help"
	"github.com/jdbann/forestry/model/stack"
	"github.com/jdbann/forestry/model/world"
	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/geo"
)

type Model struct {
	Keys  help.KeyMap
	Stack tea.Model
}

type Params struct {
	Client *client.Client
}

func New(params Params) Model {
	keys := help.DefaultKeys
	helpModel := help.New(help.Params{
		Keys: keys,
	})
	worldModel := world.New(world.Params{
		Client:  params.Client,
		MapSize: geo.Size{Width: 64, Height: 24},
	})
	stackModel := stack.NewVertical(
		stack.FlexSlot(worldModel),
		stack.FixedSlot(helpModel),
	)

	return Model{
		Keys:  keys,
		Stack: stackModel,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Stack.Init(), world.AddPerson, world.AddTree, world.AddTree, world.AddTree, world.AddTree, world.AddTree)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.Keys.AddPerson):
			return m, world.AddPerson
		case key.Matches(msg, m.Keys.AddTree):
			return m, world.AddTree
		default:
		}
	default:
	}

	var cmd tea.Cmd
	m.Stack, cmd = m.Stack.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.Stack.View()
}
