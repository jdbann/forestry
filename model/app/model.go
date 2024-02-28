package app

import (
	"math/rand"

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
	Rng    *rand.Rand
}

func New(params Params) Model {
	keys := help.DefaultKeys
	helpModel := help.New(help.Params{
		Keys: keys,
	})
	worldModel := world.New(world.Params{
		Client:  params.Client,
		Rng:     params.Rng,
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
	return tea.Batch(m.Stack.Init(), world.AddPerson)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			return m, tea.Quit
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
