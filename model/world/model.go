package world

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/geo"
)

var mapStyle = lipgloss.NewStyle().Background(color.Grass3).Foreground(color.Grass11)

type Model struct {
	Entities []Entity
	Size     geo.Size
}

type Entity interface {
	Position() geo.Point
	View() string
}

func New(size geo.Size) Model {
	return Model{
		Entities: []Entity{},
		Size:     size,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AddEntityMsg:
		m.Entities = append(m.Entities, msg)
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
			for _, entity := range m.Entities {
				if entity.Position().Equals(geo.Point{X: x, Y: y}) {
					b.WriteString(entity.View())
					continue PixelLoop
				}
			}

			b.WriteString(" ")
		}
	}

	return mapStyle.Render(b.String())
}

type AddEntityMsg Entity

func AddEntity(e Entity) tea.Msg {
	return AddEntityMsg(e)
}
