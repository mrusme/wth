package tui

import (
	"errors"
	"fmt"
	"time"

	"github.com/76creates/stickers"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	// lib "github.com/mrusme/libwth"
	"github.com/mrusme/libwth/cfg"
	"github.com/mrusme/libwth/module"
	"go.uber.org/zap"
)

type KeyMap struct {
	Up   key.Binding
	Down key.Binding
	Quit key.Binding
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

type ModuleMeta struct {
	Width  int
	Height int
}

type Model struct {
	keymap        KeyMap
	config        *cfg.Cfg
	modules       *[]*module.Module
	meta          map[string]ModuleMeta
	moduleUpdates []time.Time
	flexbox       *stickers.FlexBox
	screen        []int
	currentFocus  int
	pingChan      chan module.HeartbeatMsg
	log           *zap.SugaredLogger
}

func New(config *cfg.Cfg, modules *[]*module.Module, log *zap.SugaredLogger) Model {
	m := Model{
		keymap:       DefaultKeyMap,
		config:       config,
		modules:      modules,
		meta:         make(map[string]ModuleMeta),
		flexbox:      stickers.NewFlexBox(0, 0),
		screen:       []int{0, 0},
		currentFocus: -1,
		pingChan:     make(chan module.HeartbeatMsg),
		log:          log,
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

func (m Model) ping() tea.Cmd {
	return func() tea.Msg {
		for {
			m.pingChan <- module.HeartbeatMsg{
				Now: time.Now(),
			}
			time.Sleep(time.Second)
		}
	}
}

func (m Model) pong() tea.Cmd {
	return func() tea.Msg {
		return module.HeartbeatMsg(<-m.pingChan)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.ping(),
		m.pong(),
	)
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
			id := m.config.Modules[i].ID
			cell := m.getCellByModuleID(id)
			if cell != nil {
				m.meta[id] = ModuleMeta{
					Width:  cell.GetWidth(),
					Height: cell.GetHeight(),
				}
				_msg := module.ModuleResizeEvent{
					Width:  cell.GetWidth(),
					Height: cell.GetHeight(),
				}
				v, cmd := (*(*m.modules)[i]).Update(_msg)
				(*(*m.modules)[i]) = v
				cmds = append(cmds, cmd)
			}
		}

	case module.HeartbeatMsg:
		for i := 0; i < len(*m.modules); i++ {
			if len(m.moduleUpdates) <= i {
				m.moduleUpdates = append(m.moduleUpdates, msg.Now)
			}

			refreshInterval, err := time.ParseDuration(m.config.Modules[i].RefreshInterval)
			if err != nil {
				continue
			}

			nextUpdate := m.moduleUpdates[i].Add(refreshInterval)

			if msg.IsDueNow(nextUpdate) {
				m.moduleUpdates[i] = msg.Now

				v, cmd := (*(*m.modules)[i]).Update(msg)
				(*(*m.modules)[i]) = v
				cmds = append(cmds, cmd)
			}
		}
		cmds = append(cmds, m.pong())

	}

	/*
	  for i := 0; i < len(*m.modules); i++ {
	    v, cmd := (*(*m.modules)[i]).Update(msg)
	    (*(*m.modules)[i]) = v
	    cmds = append(cmds, cmd)
	  }
	*/

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.flexbox.ForceRecalculate()

	for rowIdx, row := range m.config.Layout.Rows {
		for cellIdx, cell := range row.Cells {
			var content string

			if cell.ModuleID == "-" {
				content = ""
			} else {
				for i := 0; i < len(*m.modules); i++ {
					id := m.config.Modules[i].ID
					if id == cell.ModuleID {
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

func (m Model) getCellByModuleID(moduleID string) *stickers.FlexBoxCell {
	for rowIdx, row := range m.config.Layout.Rows {
		for cellIdx, cell := range row.Cells {
			if cell.ModuleID == moduleID {
				flexRow := m.flexbox.Row(rowIdx)
				if flexRow == nil {
					return nil
				}
				flexCell := flexRow.Cell(cellIdx)
				if flexCell == nil {
					return nil
				}

				return flexCell
			}
		}
	}

	return nil
}

func (m Model) setSizes(winWidth int, winHeight int) {
	m.screen[0] = winWidth
	m.screen[1] = winHeight
	m.flexbox.SetWidth(winWidth)
	m.flexbox.SetHeight(winHeight)
	m.flexbox.ForceRecalculate()
}
