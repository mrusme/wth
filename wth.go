package main

import (
  "plugin"

	lib "github.com/mrusme/libwth"
  "github.com/mrusme/wth/tui"
  tea "github.com/charmbracelet/bubbletea"
)

type Runtime struct {
  Config                 lib.Cfg
  Modules                []*lib.Module
}

func main() {
  rt := new(Runtime)

  config, err := lib.NewCfg()
  if err != nil {
    panic(err)
  }
  rt.Config = config

  for _, cfgModule := range config.Modules {
    ctx, err := lib.NewCtx(&config, cfgModule.ID)
    if err != nil {
      panic(err)
    }

    mod, err := plugin.Open(cfgModule.Path)
    if err != nil {
      panic(err)
    }

    symNewModule, err := mod.Lookup("NewModule")
    if err != nil {
      panic(err)
    }

    module, err := symNewModule.(func(*lib.Ctx) (lib.Module, error))(ctx)
    if err != nil {
      panic(err)
    }
    rt.Modules = append(rt.Modules, &module)
  }

  t := tea.NewProgram(tui.New(&rt.Config, &rt.Modules), tea.WithAltScreen())
  err = t.Start()
  if err != nil {
    panic(err)
  }
}

