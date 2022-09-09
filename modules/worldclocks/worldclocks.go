package main

import (
  "fmt"
  "strings"
  "time"

  lib "github.com/mrusme/libwth"
  "github.com/mrusme/libwth/module"

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

  locations       []*time.Location
  locationStyle   lipgloss.Style

  timeformat      string
  dateStyle       lipgloss.Style

}

func NewModule(ctx *lib.Ctx) (module.Module, error) {
  module := new(Module)
  module.ctx = ctx

  tzlist := strings.Split(module.ctx.ConfigValue("timezones"), ";")
  for _, tz := range tzlist {
    loc, err := time.LoadLocation(tz)
    if err != nil {
      module.ctx.Log.Error(err)
      continue
    }
    module.locations = append(module.locations, loc)
  }
  module.locationStyle = ctx.Theme().DefaultLabelStyle()

  module.timeformat = "15:04:05"
  if module.ctx.Module.RefreshInterval != "1s" {
    module.timeformat = "15:04"
  }

  module.dateStyle = ctx.Theme().DefaultTextMutedStyle()
  module.viewportStyle = ctx.Theme().DefaultModuleViewStyle()

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

  case module.ModuleResizeEvent:
    m.viewportStyle.Width(msg.Width - 4)
    m.viewportStyle.Height(msg.Height - 4)
    m.viewport = viewport.New(msg.Width - 4, msg.Height - 4)
    m.viewport.Width =  msg.Width - 4
    m.viewport.Height = msg.Height - 4

  }

  var cmd tea.Cmd

  var content string = ""
  for _, location := range m.locations {
    content = fmt.Sprintf(
      "%s%s: %*s %s\n",
      content,
      m.locationStyle.Render(location.String()),
      (30-len(location.String())),
      time.Now().In(location).Format(m.timeformat),
      m.dateStyle.Render(time.Now().In(location).Format("Jan 2")),
    )
  }
  m.viewport.SetContent(content)
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

