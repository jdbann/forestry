package world

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/component/agent"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

var frameRate = time.Second / 60
var mapStyle = lipgloss.NewStyle().Background(color.Grass3).Foreground(color.Grass11)

type Model struct {
	RenderSystem *render.System
	Rng          *rand.Rand
	Scene        *ecs.Scene
	Size         geo.Size
}

type Params struct {
	Rng  *rand.Rand
	Size geo.Size
}

func New(params Params) Model {
	scene := &ecs.Scene{}
	renderSystem := &render.System{}
	scene.AddSystem(renderSystem)
	scene.AddSystem(&agent.System{WorldSize: params.Size, Rng: params.Rng})

	return Model{
		RenderSystem: renderSystem,
		Rng:          params.Rng,
		Scene:        scene,
		Size:         params.Size,
	}
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AddEntityMsg:
		m.Scene.AddEntity(msg)

	case AddPersonMsg:
		m.Scene.AddEntity(newPerson(m.Rng, m.Size))

	case TickMsg:
		m.Scene.Update(time.Duration(msg))
		return m, doTick()
	}

	return m, nil
}

func (m Model) View() string {
	b := &strings.Builder{}

	// This is an inefficient method for drawing the map but it is good enough
	// for now.
	for y := 0; y < m.Size.Height; y++ {
		if y > 0 {
			b.WriteString("\n")
		}

	PixelLoop:
		for x := 0; x < m.Size.Width; x++ {
			for _, component := range m.RenderSystem.Components {
				if component.Position.Equals(geo.Point{X: x, Y: y}) {
					b.WriteString(component.View())
					continue PixelLoop
				}
			}

			b.WriteString(" ")
		}
	}

	return mapStyle.Render(b.String())
}

func doTick() tea.Cmd {
	tickRequested := time.Now()
	return tea.Tick(frameRate, func(t time.Time) tea.Msg {
		return TickMsg(time.Since(tickRequested))
	})
}

func newPerson(rng *rand.Rand, bounds geo.Size) *ecs.Entity {
	return ecs.NewEntity(
		&render.Component{
			Position: bounds.PointWithin(rng),
			Rune:     'P',
		},
		&agent.Component{
			StepFrequency: time.Millisecond * 500,
			SinceLastStep: 0,
		},
	)
}

type TickMsg time.Duration
type AddEntityMsg *ecs.Entity
type AddPersonMsg struct{}

func AddPerson() tea.Msg {
	return AddPersonMsg{}
}
