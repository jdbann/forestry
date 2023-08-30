package person

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdbann/forestry/component/brain"
	"github.com/jdbann/forestry/component/render"
	"github.com/jdbann/forestry/pkg/ecs"
)

func BrainComponent(id int) *brain.Component {
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
	return tea.Sequence(tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return nil
	}), attemptRegistration(c.Entity.ID()))
}

type RegisteredState struct {
	brain.BaseState
}

func (s RegisteredState) Update(_ *brain.Component, _ tea.Msg) tea.Cmd {
	fmt.Println("RegisteredState - Update")
	return nil
}

func attemptRegistration(id int) tea.Cmd {
	return func() tea.Msg {
		body := fmt.Sprintf(`{"id": %d}`, id)
		request, err := http.NewRequest(http.MethodPost, "http://localhost:3000/people", strings.NewReader(body))
		if err != nil {
			return errMsg(err)
		}

		res, err := http.DefaultClient.Do(request)
		if err != nil {
			return errMsg(err)
		}

		if res.StatusCode != 200 {
			return errMsg(fmt.Errorf("registering: %s", res.Status))
		}

		return registeredMsg{}
	}
}

type errMsg error
type registeredMsg struct{}
