package stack_test

import (
	"bytes"
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/jdbann/forestry/model/box"
	"github.com/jdbann/forestry/model/stack"
	"github.com/jdbann/forestry/pkg/color"
	"github.com/jdbann/forestry/pkg/geo"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

func TestModel(t *testing.T) {
	type testCase struct {
		name  string
		model stack.Model
	}

	run := func(t *testing.T, tc testCase) {
		tm := teatest.NewTestModel(t, tc.model, teatest.WithInitialTermSize(28, 12))
		var out bytes.Buffer
		teatest.WaitFor(t, io.TeeReader(tm.Output(), &out), func(bts []byte) bool {
			return bytes.Contains(bts, []byte("Flexi"))
		})

		runStep(t, "initial", func(t *testing.T) {
			out, err := io.ReadAll(&out)
			if err != nil {
				t.Fatal(err)
			}
			teatest.RequireEqualOutput(t, out)
		})

		tm.Send(tea.WindowSizeMsg{Width: 64, Height: 20})

		runStep(t, "resized", func(t *testing.T) {
			tm.Send(tea.Quit())
			out, err := io.ReadAll(tm.FinalOutput(t))
			if err != nil {
				t.Error(err)
			}
			teatest.RequireEqualOutput(t, out)
		})
	}

	testCases := []testCase{
		{
			name: "flexible vertical stack",
			model: stack.NewVertical(
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Gray9, Text: "Flexi"})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Gray9, Text: "Flexi"})),
			),
		},
		{
			name: "flexible horizontal stack",
			model: stack.NewHorizontal(
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Gray9, Text: "Flexi"})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Gray9, Text: "Flexi"})),
			),
		},
		{
			name: "part-fixed vertical stack",
			model: stack.NewVertical(
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
				stack.FixedSlot(box.New(box.Params{Fill: color.Gray9, Text: "Fixed", Size: geo.Size{Height: 10}})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
			),
		},
		{
			name: "part-fixed horizontal stack",
			model: stack.NewHorizontal(
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
				stack.FixedSlot(box.New(box.Params{Fill: color.Gray9, Text: "Fixed", Size: geo.Size{Width: 7}})),
				stack.FlexSlot(box.New(box.Params{Fill: color.Grass9, Text: "Flexi"})),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func runStep(t *testing.T, name string, fn func(t *testing.T)) {
	if !t.Run(name, fn) {
		t.FailNow()
	}
}
