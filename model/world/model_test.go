package world_test

import (
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/jdbann/forestry/model/world"
	"github.com/jdbann/forestry/pkg/geo"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

func TestModel(t *testing.T) {
	m := world.New(geo.Size{Width: 64, Height: 32})
	tm := teatest.NewTestModel(t, m)
	tm.Send(world.AddEntity(runeEntity{
		Pos:  geo.Point{X: 4, Y: 4},
		Rune: 'P',
	}))
	tm.Send(tea.Quit())
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

type runeEntity struct {
	Pos  geo.Point
	Rune rune
}

func (e runeEntity) Position() geo.Point {
	return e.Pos
}

func (e runeEntity) View() string {
	return string(e.Rune)
}
