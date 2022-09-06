package main

import (
  "plugin"

  "github.com/mrusme/wth/lib"
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
    ctx, err := lib.NewCtx(&config, cfgModule.Path)
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

    var newModule lib.NewModule
    newModule, ok := symNewModule.(lib.NewModule)
    if !ok {
      panic("issue with module symbol, contact the module developer")
    }

    module, err := newModule(ctx)
    if err != nil {
      panic(err)
    }
    rt.Modules = append(rt.Modules, module)
  }

  tui := tea.NewProgram(tui.New(&rt.Config, &rt.Modules), tea.WithAltScreen())
  err = tui.Start()
  if err != nil {
    panic(err)
  }
}

