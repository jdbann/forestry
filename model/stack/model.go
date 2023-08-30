package stack

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdbann/forestry/pkg/geo"
)

type Model struct {
	Distributor Distributor
	Size        geo.Size
	Slots       []Slot
}

type Params struct {
	Distributor Distributor
	Slots       []Slot
}

func NewVertical(slots ...Slot) Model {
	return New(Params{
		Distributor: HeightDistributor{},
		Slots:       slots,
	})
}

func NewHorizontal(slots ...Slot) Model {
	return New(Params{
		Distributor: WidthDistributor{},
		Slots:       slots,
	})
}

func New(params Params) Model {
	return Model{
		Distributor: params.Distributor,
		Slots:       params.Slots,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, slot := range m.Slots {
		cmds = append(cmds, slot.Model.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.Size = geo.Size(msg)

		remainingSlots, remainingSize := float64(len(m.Slots)), float64(m.Distributor.AvailableSize(msg))

		for i, slot := range m.Slots {
			if !slot.Fixed {
				continue
			}

			slotSize := m.Distributor.SlotSize(slot)
			m.Slots[i].Model, cmd = m.Distributor.UpdateSlot(msg, slot.Model, slotSize)
			remainingSlots--
			remainingSize -= float64(slotSize)
			cmds = append(cmds, cmd)
		}

		for i, slot := range m.Slots {
			if slot.Fixed {
				continue
			}

			slotSize := int(math.Round(remainingSize / remainingSlots))
			m.Slots[i].Model, cmd = m.Distributor.UpdateSlot(msg, slot.Model, slotSize)
			remainingSlots--
			remainingSize -= float64(slotSize)
			cmds = append(cmds, cmd)
		}

	default:
		for i, slot := range m.Slots {
			m.Slots[i].Model, cmd = slot.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var views []string
	for _, slot := range m.Slots {
		views = append(views, slot.Model.View())
	}
	return m.Distributor.JoinViews(views)
}

type Slot struct {
	Model tea.Model
	Fixed bool
}

func FixedSlot(m tea.Model) Slot {
	return Slot{
		Model: m,
		Fixed: true,
	}
}

func FlexSlot(m tea.Model) Slot {
	return Slot{
		Model: m,
		Fixed: false,
	}
}

type Distributor interface {
	AvailableSize(tea.WindowSizeMsg) int
	SlotSize(Slot) int
	UpdateSlot(tea.WindowSizeMsg, tea.Model, int) (tea.Model, tea.Cmd)
	JoinViews([]string) string
}

type HeightDistributor struct{}

func (d HeightDistributor) AvailableSize(msg tea.WindowSizeMsg) int {
	return msg.Height
}

func (d HeightDistributor) SlotSize(slot Slot) int {
	return lipgloss.Height(slot.Model.View())
}

func (d HeightDistributor) UpdateSlot(msg tea.WindowSizeMsg, model tea.Model, size int) (tea.Model, tea.Cmd) {
	return model.Update(tea.WindowSizeMsg{Height: size, Width: msg.Width})
}

func (d HeightDistributor) JoinViews(views []string) string {
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

type WidthDistributor struct{}

func (d WidthDistributor) AvailableSize(msg tea.WindowSizeMsg) int {
	return msg.Width
}

func (d WidthDistributor) SlotSize(slot Slot) int {
	return lipgloss.Width(slot.Model.View())
}

func (d WidthDistributor) UpdateSlot(msg tea.WindowSizeMsg, model tea.Model, size int) (tea.Model, tea.Cmd) {
	return model.Update(tea.WindowSizeMsg{Height: msg.Height, Width: size})
}

func (d WidthDistributor) JoinViews(views []string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}
