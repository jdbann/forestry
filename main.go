package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/model/app"
)

func main() {
	m := app.New(app.Params{
		Rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	})
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
