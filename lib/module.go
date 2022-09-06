package lib

import (
  tea "github.com/charmbracelet/bubbletea"
)

type Module interface {
  View() (string)
  Update(msg tea.Msg) (tea.Model, tea.Cmd)
}

type NewModule func(*Ctx) (*Module, error)

