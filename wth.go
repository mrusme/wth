package main

import (
  "plugin"
  "runtime"
  "net/url"
  "os"

	lib "github.com/mrusme/libwth"
  "github.com/mrusme/wth/tui"
  tea "github.com/charmbracelet/bubbletea"

  "go.uber.org/zap"
)

type Runtime struct {
  Config                 lib.Cfg
  Modules                []*lib.Module
  Logger                 *zap.Logger
  Log                    *zap.SugaredLogger
}

func NewLogger(filename string) (*zap.Logger, error) {
  if runtime.GOOS == "windows" {
    zap.RegisterSink("winfile", func(u *url.URL) (zap.Sink, error) {
      return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    })
  }

  cfg := zap.NewDevelopmentConfig()
  if runtime.GOOS == "windows" {
    cfg.OutputPaths = []string{
      "stdout",
      "winfile:///" + filename,
    }
  } else {
    cfg.OutputPaths = []string{
      filename,
    }
  }

  return cfg.Build()
}


func main() {
  rt := new(Runtime)

  config, err := lib.NewCfg()
  if err != nil {
    panic(err)
  }
  rt.Config = config

  logger, err := NewLogger("wth.log")
  if err != nil {
    panic(err)
  }
  rt.Logger = logger
  defer rt.Logger.Sync()
  rt.Log = rt.Logger.Sugar()

  for _, cfgModule := range config.Modules {
    ctx, err := lib.NewCtx(
      &config,
      cfgModule.ID,
    )
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

  t := tea.NewProgram(tui.New(&rt.Config, &rt.Modules, rt.Log), tea.WithAltScreen())
  err = t.Start()
  if err != nil {
    panic(err)
  }
}

