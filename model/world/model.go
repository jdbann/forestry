package world

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

var mapStyle = lipgloss.NewStyle().Background(color.Grass3).Foreground(color.Grass11)

type Model struct {
	RenderSystem *render.System
	Scene        *ecs.Scene
	Size         geo.Size
}

func New(size geo.Size) Model {
	scene := &ecs.Scene{}
	renderSystem := &render.System{}
	scene.AddSystem(renderSystem)

	return Model{
		RenderSystem: renderSystem,
		Scene:        scene,
		Size:         size,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AddEntityMsg:
		m.Scene.AddEntity(msg)
		return m, nil
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

type AddEntityMsg *ecs.Entity
