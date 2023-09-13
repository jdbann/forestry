package app_test

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"

	"github.com/jdbann/forestry/model/app"
	"github.com/jdbann/forestry/pkg/client"
)

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

func TestModel(t *testing.T) {
	c, _ := client.New("http://local.test:3000")
	m := app.New(app.Params{
		Client: c,
		Rng:    rand.New(rand.NewSource(98)),
	})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(72, 32))
	var out bytes.Buffer
	teatest.WaitFor(t, io.TeeReader(tm.Output(), &out), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("quit"))
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	teatest.RequireEqualOutput(t, out.Bytes())
}
