package world_test

import (
	"io"
	"math/rand"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/jdbann/forestry/model/world"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

func TestModel(t *testing.T) {
	m := world.New(world.Params{
		Rng:     rand.New(rand.NewSource(98)),
		MapSize: geo.Size{Width: 64, Height: 32},
	})
	tm := teatest.NewTestModel(t, m)
	tm.Send(world.AddPersonMsg{})
	tm.Send(world.TickMsg(time.Second + 1))
	tm.Send(tea.Quit())
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func newEntity(components ...ecs.Component) *ecs.Entity {
	e := &ecs.Entity{}
	for _, component := range components {
		ecs.AddComponent(e, component)
	}
	return e
}
