package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdbann/forestry/model/app"
	"github.com/jdbann/forestry/pkg/client"
)

var baseURL = flag.String("base-url", "http://localhost:3000", "base URL for your API")

func main() {
	flag.Parse()

	c, err := client.New(*baseURL)
	if err != nil {
		_, _ = fmt.Println(err)
		os.Exit(1)
	}

	m := app.New(app.Params{
		Client: c,
		Rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	})
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		_, _ = fmt.Println(err)
		os.Exit(1)
	}
}
