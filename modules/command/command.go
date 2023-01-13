package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
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
	envs    []string
	cmd     *exec.Cmd
	width   int
	height  int
}

func NewModule(ctx *lib.Ctx) (module.Module, error) {
	module := new(Module)
	module.ctx = ctx
	module.envs = os.Environ()

	command := module.ctx.ConfigValue("command")
	if command == "" {
		command = "echo \"No command specified\""
	}

	var err error
	var envs []string
	envs, module.command, err = shellwords.ParseWithEnvs(command)
	if err != nil {
		return nil, err
	}
	module.envs = append(module.envs, envs...)

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

	}

	var cmd tea.Cmd

	var content string = ""
	m.cmd = exec.Command(m.command[0], m.command[1:]...)
	m.cmd.Env = m.envs
	m.cmd.Env = append(m.cmd.Env,
		fmt.Sprintf("COLUMNS=%d", m.width-4),
		fmt.Sprintf("HEIGHT=%d", m.height-4),
	)
	out, err := m.cmd.Output()
	if err != nil {
		content = fmt.Sprintf("%s %v:\n%s\n", m.command[0], m.command, err.Error())
	} else {
		content = string(out)
	}

	contentLines := strings.Split(content, "\n")
	startIdx := 0
	endIdx := len(contentLines) - 1
	if endIdx > m.height {
		startIdx = endIdx - m.height + 4
	}
	content = ""
	for i := startIdx; i <= endIdx; i++ {
		content += exstrings.SubString(contentLines[i], 0, m.width-6) + "\n"
	}

	m.viewport.SetContent(content)
	m.viewport, cmd = m.viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Module) View() string {
	return m.viewportStyle.Render(m.viewport.View())
}
