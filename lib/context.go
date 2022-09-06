package lib

import (
  "errors"
  "fmt"
)

type Ctx struct {
  config       *Cfg
  ModulePath   string
  Module       *CfgModule
}

func NewCtx(config *Cfg, path string) (*Ctx, error) {
  ctx := new(Ctx)
  ctx.config = config
  ctx.ModulePath = path
  ctx.Module = nil
  for i := 0; i < len(ctx.config.Modules); i++ {
    if ctx.config.Modules[i].Path == path {
      ctx.Module = &ctx.config.Modules[i]
      break
    }
  }
  if ctx.Module == nil {
    return nil, errors.New(
      fmt.Sprint(
        "No module configuration with path %s available!",
        path,
      ),
    )
  }
  ctx.ModulePath = path

  return ctx, nil
}

