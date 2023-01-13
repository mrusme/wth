package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	lib "github.com/mrusme/libwth"
	"github.com/mrusme/libwth/module"
	"github.com/thinkeridea/go-extend/exstrings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Module struct {
	ctx           *lib.Ctx
	viewport      viewport.Model
	viewportStyle lipgloss.Style

	locations     []*time.Location
	locationStyle lipgloss.Style

	command []string
	cmd     *exec.Cmd
	width   int
	height  int
}

func NewModule(ctx *lib.Ctx) (module.Module, error) {
	module := new(Module)
	module.ctx = ctx

	command := module.ctx.ConfigValue("command")
	if command == "" {
		command = "echo No command specified"
	}
	module.command = strings.SplitN(command, " ", 2)
	module.cmd = exec.Command(module.command[0], module.command[1])

	module.viewportStyle = ctx.Theme().DefaultModuleViewStyle()

	return module, nil
}

func (m Module) Init() tea.Cmd {
	return nil
}

func (m Module) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case module.ModuleResizeEvent:
		m.width = msg.Width
		m.height = msg.Height
		m.viewportStyle.Width(msg.Width - 4)
		m.viewportStyle.Height(msg.Height - 4)
		m.viewport = viewport.New(msg.Width-4, msg.Height-4)
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 4
		m.cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("COLUMNS=%d", msg.Width-4),
			fmt.Sprintf("HEIGHT=%d", msg.Height-4),
		)

	}

	var cmd tea.Cmd

	var content string = ""
	out, err := m.cmd.Output()
	if err != nil {
		content = err.Error()
	} else {
		content = string(out)
	}

	contentLines := strings.Split(content, "\n")
	for i := 0; i < m.height; i++ {
		contentLines[i] = exstrings.SubString(contentLines[i], 0, m.width-6)
	}
	if len(contentLines) > m.height {
		content = strings.Join(contentLines[0:m.height], "\n")
	}

	m.viewport.SetContent(content)
	m.viewport, cmd = m.viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Module) View() string {
	return m.viewportStyle.Render(m.viewport.View())
}
