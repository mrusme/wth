package tui

import (
	// "fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lib "github.com/mrusme/libwth"
)

type KeyMap struct {
    FirstTab      key.Binding
    SecondTab     key.Binding
    ThirdTab      key.Binding
    FourthTab     key.Binding
    FifthTab      key.Binding
    SixthTab      key.Binding
    SeventhTab    key.Binding
    EightTab      key.Binding
    NinthTab      key.Binding
    TenthTab      key.Binding
    EleventhTab   key.Binding
    TwelfthTab    key.Binding
    ThirteenthTab key.Binding
    PrevTab       key.Binding
    NextTab       key.Binding
    Up            key.Binding
    Down          key.Binding
    Quit          key.Binding
}

var DefaultKeyMap = KeyMap{
  FirstTab: key.NewBinding(
    key.WithKeys("f1"),
    key.WithHelp("f1", "first tab"),
  ),
  SecondTab: key.NewBinding(
    key.WithKeys("f2"),
    key.WithHelp("f2", "second tab"),
  ),
  ThirdTab: key.NewBinding(
    key.WithKeys("f3"),
    key.WithHelp("f3", "third tab"),
  ),
  FourthTab: key.NewBinding(
    key.WithKeys("f4"),
    key.WithHelp("f4", "fourth tab"),
  ),
  FifthTab: key.NewBinding(
    key.WithKeys("f5"),
    key.WithHelp("f5", "fifth tab"),
  ),
  SixthTab: key.NewBinding(
    key.WithKeys("f6"),
    key.WithHelp("f6", "sixth tab"),
  ),
  SeventhTab: key.NewBinding(
    key.WithKeys("f7"),
    key.WithHelp("f7", "seventh tab"),
  ),
  EightTab: key.NewBinding(
    key.WithKeys("f8"),
    key.WithHelp("f8", "eight tab"),
  ),
  NinthTab: key.NewBinding(
    key.WithKeys("f9"),
    key.WithHelp("f9", "ninth tab"),
  ),
  TenthTab: key.NewBinding(
    key.WithKeys("f10"),
    key.WithHelp("f10", "tenth tab"),
  ),
  EleventhTab: key.NewBinding(
    key.WithKeys("f11"),
    key.WithHelp("f11", "eleventh tab"),
  ),
  TwelfthTab: key.NewBinding(
    key.WithKeys("f12"),
    key.WithHelp("f12", "twelfth tab"),
  ),
  ThirteenthTab: key.NewBinding(
    key.WithKeys("f13"),
    key.WithHelp("f13", "thirteenth tab"),
  ),
  PrevTab: key.NewBinding(
    key.WithKeys("ctrl+p"),
    key.WithHelp("ctrl+p", "previous tab"),
  ),
  NextTab: key.NewBinding(
    key.WithKeys("ctrl+n"),
    key.WithHelp("ctrl+n", "next tab"),
  ),
  Up: key.NewBinding(
    key.WithKeys("k", "up"),
    key.WithHelp("↑/k", "move up"),
  ),
  Down: key.NewBinding(
    key.WithKeys("j", "down"),
    key.WithHelp("↓/j", "move down"),
  ),
  Quit: key.NewBinding(
    key.WithKeys("q", "ctrl+q"),
    key.WithHelp("q/Q", "quit"),
  ),
}


type Model struct {
  keymap        KeyMap
  config        *lib.Cfg
  modules       *[]*lib.Module
  screen        []int
  currentFocus  int
}

func New(config *lib.Cfg, modules *[]*lib.Module) Model {
  m := Model{
    keymap:        DefaultKeyMap,
    config:        config,
    modules:       modules,
    screen:        []int{0, 0},
    currentFocus:  -1,
  }

  return m
}

func (m Model) Init() tea.Cmd {
  return tea.Batch(tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  cmds := make([]tea.Cmd, 0)

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keymap.Quit):
      return m, tea.Quit
    }

  case tea.WindowSizeMsg:
    m.setSizes(msg.Width, msg.Height)
    for i := range *m.modules {
      v, cmd := (*(*m.modules)[i]).Update(msg)
      (*(*m.modules)[i]) = v
      cmds = append(cmds, cmd)
    }
  }

  for i := 0; i < len(*m.modules); i++ {
    v, cmd := (*(*m.modules)[i]).Update(msg)
    (*(*m.modules)[i]) = v
    cmds = append(cmds, cmd)
  }

  return m, tea.Batch(cmds...)
}

func (m Model) View() (string) {
  s := strings.Builder{}

  for i := 0; i < len(*m.modules); i++ {
    s.WriteString((*(*m.modules)[i]).View())
  }
  return s.String()
}

func (m Model) setSizes(winWidth int, winHeight int) {
  m.screen[0] = winWidth
  m.screen[1] = winHeight
}

