package main

import (
  lib "github.com/mrusme/libwth"

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

type Module struct {
  keymap          KeyMap
  ctx             *lib.Ctx
  viewport        viewport.Model
  viewportStyle   lipgloss.Style
}

func NewModule(ctx *lib.Ctx) (lib.Module, error) {
  module := new(Module)
  module.ctx = ctx

  module.viewportStyle = lib.DefaultModuleViewStyle(ctx)

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

  case lib.ModuleResizeEvent:
    m.viewportStyle.Width(msg.Width - 4)
    m.viewportStyle.Height(msg.Height - 4)
    m.viewport = viewport.New(msg.Width - 4, msg.Height - 4)
    m.viewport.Width =  msg.Width - 4
    m.viewport.Height = msg.Height - 4

  }

  var cmd tea.Cmd

  m.viewport.SetContent("Hello World!")
  m.viewport, cmd = m.viewport.Update(msg)

  cmds = append(cmds, cmd)
  return m, tea.Batch(cmds...)
}

func (m Module) View() (string) {
  return m.viewportStyle.Render(m.viewport.View())
}

func (m *Module) refresh() (tea.Cmd) {
  return func () (tea.Msg) {
    // TODO: Refresh things
    return nil
  }
}

