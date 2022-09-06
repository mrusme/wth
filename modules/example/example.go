package main

import (
  "github.com/mrusme/wth/lib"

  "github.com/charmbracelet/bubbles/key"
  "github.com/charmbracelet/bubbles/viewport"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
    Refresh       key.Binding
    Select        key.Binding
    GoBack        key.Binding
    SwitchFocus   key.Binding
}

var DefaultKeyMap = KeyMap{
  Refresh: key.NewBinding(
    key.WithKeys("r", "R"),
    key.WithHelp("r/R", "refresh"),
  ),
}

var viewportStyle = lipgloss.NewStyle().
    Margin(0, 0, 0, 0).
    Padding(1, 1).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD")).
    BorderTop(true).
    BorderLeft(true).
    BorderRight(true).
    BorderBottom(true)

type Module struct {
  keymap          KeyMap
  ctx             *lib.Ctx
  viewport        viewport.Model
}

func NewModule(ctx *lib.Ctx) (*Module, error) {
  module := new(Module)
  module.ctx = ctx

  return module, nil
}

func (m Module) Init() tea.Cmd {
  return nil
}

func (m Module) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmds []tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keymap.Refresh):
      cmds = append(cmds, m.refresh())
    }
  }

  var cmd tea.Cmd

  m.viewport.SetContent("Hello World!")
  m.viewport, cmd = m.viewport.Update(msg)

  cmds = append(cmds, cmd)
  return m, tea.Batch(cmds...)
}

func (m Module) View() (string) {
  return viewportStyle.Render(m.viewport.View())
}

func (m *Module) refresh() (tea.Cmd) {
  return func () (tea.Msg) {
    // TODO: Refresh things
    return nil
  }
}

