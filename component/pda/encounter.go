package pda

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/ecs"
	"github.com/jdbann/forestry/pkg/geo"
)

type discoveryRequest struct {
	X, Y int
}

type discoveryFailMsg struct {
	err error
}

type discoverySuccessMsg struct {
	id string
}

func reportDiscovery(c *Component, position geo.Point) tea.Cmd {
	return func() tea.Msg {
		req := discoveryRequest{
			X: position.X,
			Y: position.Y,
		}
		code, response, err := client.MakeRequest[registerPayload](c.Client, http.MethodPost, "/discovery", req)
		if err != nil {
			return ecs.EntityMsg{
				EntityID: c.Entity.ID(),
				Msg:      discoveryFailMsg{err},
			}
		}

		if code != http.StatusOK {
			return ecs.EntityMsg{
				EntityID: c.Entity.ID(),
				Msg:      discoveryFailMsg{fmt.Errorf("unexpected response code: %d", code)},
			}
		}

		return ecs.EntityMsg{
			EntityID: c.Entity.ID(),
			Msg: discoverySuccessMsg{
				id: response.ID,
			},
		}
	}
}
