package person

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdbann/forestry/component/brain"
	"github.com/jdbann/forestry/component/pda"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/client"
	"github.com/jdbann/forestry/pkg/ecs"
)

func BrainComponent() *brain.Component {
	unregistered := &UnregisteredState{}
	return &brain.Component{
		AvailableStates: []brain.State{unregistered, &RegisteredState{}},
		CurrentState:    unregistered,
	}
}

type UnregisteredState struct {
	err error
}

func (s *UnregisteredState) Update(c *brain.Component, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case errMsg:
		rc, ok := ecs.GetComponent[*render.Component](c.Entity)
		if ok {
			rc.Rune = '!'
		}

		s.err = msg
	case registeredMsg:
		rc, ok := ecs.GetComponent[*render.Component](c.Entity)
		if ok {
			rc.Rune = 'P'
		}

		s.err = nil
	}
	return nil
}

func (s UnregisteredState) OnEnter(c *brain.Component) tea.Cmd {
	pdac, ok := ecs.GetComponent[*pda.Component](c.Entity)
	if !ok {
		return nil
	}

	return tea.Sequence(tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return nil
	}), attemptRegistration(pdac.Client, c.Entity.ID()))
}

type RegisteredState struct {
	brain.BaseState
}

func (s RegisteredState) Update(_ *brain.Component, _ tea.Msg) tea.Cmd {
	_, _ = fmt.Println("RegisteredState - Update")
	return nil
}

type registerPayload struct {
	ID int `json:"id"`
}

func attemptRegistration(c *client.Client, id int) tea.Cmd {
	return func() tea.Msg {
		status, _, err := client.MakeRequest[struct{}](c, http.MethodPost, "/people", registerPayload{ID: id})
		if err != nil {
			return errMsg(err)
		}

		if status != http.StatusOK {
			return errMsg(fmt.Errorf("registering: %d", status))
		}

		return registeredMsg{}
	}
}

type (
	errMsg        error
	registeredMsg struct{}
)
