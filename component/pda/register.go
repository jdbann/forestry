package pda

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/ecs"
)

type attemptRegistrationMsg struct{}

func AttemptRegistration(e *ecs.Entity) tea.Cmd {
	return func() tea.Msg {
		return ecs.EntityMsg{
			EntityID: e.ID(),
			Msg:      attemptRegistrationMsg{},
		}
	}
}

type registerPayload struct {
	ID string `json:"id"`
}

type RegisterFailMsg struct {
	err error
}

type RegisterSuccessMsg struct {
	id string
}

func performRegistration(c *Component) tea.Cmd {
	return func() tea.Msg {
		code, response, err := client.MakeRequest[registerPayload](c.Client, http.MethodPost, "/people", nil)
		if err != nil {
			return ecs.EntityMsg{
				EntityID: c.Entity.ID(),
				Msg:      RegisterFailMsg{err},
			}
		}

		if code != http.StatusOK {
			return ecs.EntityMsg{
				EntityID: c.Entity.ID(),
				Msg:      RegisterFailMsg{fmt.Errorf("unexpected response code: %d", code)},
			}
		}

		return ecs.EntityMsg{
			EntityID: c.Entity.ID(),
			Msg: RegisterSuccessMsg{
				id: response.ID,
			},
		}
	}
}
