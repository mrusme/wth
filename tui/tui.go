package tui

import (
	"errors"
	"fmt"

	"github.com/76creates/stickers"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lib "github.com/mrusme/libwth"
)

type KeyMap struct {
    Up            key.Binding
    Down          key.Binding
    Quit          key.Binding
}

var DefaultKeyMap = KeyMap{
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
  flexbox       *stickers.FlexBox
  screen        []int
  currentFocus  int
}

func New(config *lib.Cfg, modules *[]*lib.Module) Model {
  m := Model{
    keymap:        DefaultKeyMap,
    config:        config,
    modules:       modules,
    flexbox:       stickers.NewFlexBox(0, 0),
    screen:        []int{0, 0},
    currentFocus:  -1,
  }

  var flexRows []*stickers.FlexBoxRow

  for _, row := range m.config.Layout.Rows {
    flexRow := new(stickers.FlexBoxRow)
    var flexCells []*stickers.FlexBoxCell

    for _, cell := range row.Cells {
      flexCell := stickers.
        NewFlexBoxCell(cell.RatioX, cell.RatioY)

      flexCells = append(flexCells, flexCell)
    }
    flexRow.AddCells(flexCells)
    flexRows = append(flexRows, flexRow)
  }
  m.flexbox.AddRows(flexRows)

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
	m.flexbox.ForceRecalculate()

  for rowIdx, row := range m.config.Layout.Rows {
    for cellIdx, cell := range row.Cells {
      var content string

      if cell.ModuleID == "-" {
        content = ""
      } else {
        for i := 0; i < len(*m.modules); i++ {
          if m.config.Modules[i].ID == cell.ModuleID {
            content = (*(*m.modules)[i]).View()
          }
        }
      }

      flexRow := m.flexbox.Row(rowIdx)
      if flexRow == nil {
        panic(errors.New(
          fmt.Sprintf("flex row %d could not be found", rowIdx),
        ))
      }
      flexCell := flexRow.Cell(cellIdx)
      if flexCell == nil {
        panic(errors.New(
          fmt.Sprintf("flex cell %d could not be found", cellIdx),
        ))
      }

      flexCell.SetContent(content)
    }
  }

  return m.flexbox.Render()
}

func (m Model) setSizes(winWidth int, winHeight int) {
  m.screen[0] = winWidth
  m.screen[1] = winHeight
  m.flexbox.SetWidth(winWidth)
  m.flexbox.SetHeight(winHeight)
}

