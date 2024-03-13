package world

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdbann/forestry/component/agent"
	"github.com/jdbann/forestry/component/employment"
	"github.com/jdbann/forestry/component/graph"
	"github.com/jdbann/forestry/component/pda"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/entity/person"
	"github.com/jdbann/forestry/entity/tree"
	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

var (
	frameRate = time.Second / 60
	mapStyle  = lipgloss.NewStyle().Background(color.Grass3).Foreground(color.Grass11)
)

type Model struct {
	MapSize      geo.Size
	RenderSystem *render.System
	Scene        *ecs.Scene
	Size         geo.Size
}

type Params struct {
	Client  *client.Client
	MapSize geo.Size
}

func New(params Params) Model {
	scene := &ecs.Scene{}
	renderSystem := &render.System{}
	graphSystem := graph.NewSystem(params.MapSize)

	scene.AddSystem(graphSystem)
	scene.AddSystem(&pda.System{Client: params.Client})
	scene.AddSystem(renderSystem)
	scene.AddSystem(&agent.System{WorldSize: params.MapSize})
	scene.AddSystem(&employment.System{})

	return Model{
		RenderSystem: renderSystem,
		Scene:        scene,
		MapSize:      params.MapSize,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(doTick(), m.Scene.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	switch msg := msg.(type) {
	case AddPersonMsg:
		cmd = m.Scene.AddEntity(person.New(m.MapSize.RandomPointWithin()))
		return m, cmd

	case AddTreeMsg:
		cmd = m.Scene.AddEntity(tree.New(m.MapSize.RandomPointWithin()))
		return m, cmd

	case ecs.TickMsg:
		cmds = append(cmds, doTick())

	case tea.WindowSizeMsg:
		m.Size = geo.Size(msg)
		return m, nil
	}

	cmd = m.Scene.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	b := &strings.Builder{}

	// This is an inefficient method for drawing the map but it is good enough
	// for now.
	for y := 0; y < m.MapSize.Height; y++ {
		if y > 0 {
			_, _ = b.WriteString("\n")
		}

	PixelLoop:
		for x := 0; x < m.MapSize.Width; x++ {
			for _, component := range m.RenderSystem.Components {
				if component.Position.Equals(geo.Point{X: x, Y: y}) {
					_, _ = b.WriteString(component.View())
					continue PixelLoop
				}
			}

			_, _ = b.WriteString(" ")
		}
	}

	return lipgloss.Place(
		m.Size.Width, m.Size.Height,
		lipgloss.Center, lipgloss.Center,
		mapStyle.Render(b.String()),
	)
}

func doTick() tea.Cmd {
	tickRequested := time.Now()
	return tea.Tick(frameRate, func(t time.Time) tea.Msg {
		return ecs.TickMsg(time.Since(tickRequested))
	})
}

type AddPersonMsg struct{}

func AddPerson() tea.Msg {
	return AddPersonMsg{}
}

type AddTreeMsg struct{}

func AddTree() tea.Msg {
	return AddTreeMsg{}
}
