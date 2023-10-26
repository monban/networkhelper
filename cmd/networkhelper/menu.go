package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/erikgeiser/promptkit/textinput"
	netmask "github.com/monban/bubble-netmask"
)

type Menu struct {
	index    int
	controls []tea.Model
}

type BlurFocuser interface {
	Blur()
	Focus() tea.Cmd
	Focused() bool
}

func New() Menu {
	var m Menu
	m.addControl(textinput.NewModel(textinput.New("hello")))
	m.addControl(textinput.NewModel(textinput.New("world")))
	m.addControl(netmask.New("192.168.1.1"))
	return m
}

func (m *Menu) addControl(c tea.Model) {
	m.controls = append(m.controls, c)
}

func (m Menu) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range m.controls {
		cmd := c.Init()
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, m.UpdateFocus())
	return tea.Batch(cmds...)
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab, tea.KeyEnter:
			m.index = (m.index + 1) % len(m.controls)
			cmds = append(cmds, m.UpdateFocus())
		default:
			m.controls[m.index], cmd = m.controls[m.index].Update(msg)
			cmds = append(cmds, cmd)
		}
	default:
		m.controls[0], cmd = m.controls[0].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Menu) UpdateFocus() tea.Cmd {
	var cmd tea.Cmd
	for i, c := range m.controls {
		if c, ok := c.(BlurFocuser); ok {
			if i == m.index {
				cmd = c.Focus()
			} else {
				c.Blur()
			}
		}
	}
	if cmd != nil {
		log.Info("", "cmd", cmd())
	}
	return cmd
}

func (m Menu) View() string {
	var views []string
	for _, c := range m.controls {
		views = append(views, c.View())
	}
	return lipgloss.JoinVertical(lipgloss.Center, views...)
}
