package lib

import (
  "strings"

  "github.com/spf13/viper"
)

type CfgModule struct {
  Enabled                bool
  Name                   string
  Path                   string
  Config                 map[string]string
  Position               struct {
    Top                  int
    Left                 int
    Height               int
    Width                int
  }
  RefreshInterval        string
}

type Cfg struct {
  Modules                []CfgModule
  Theme                  struct {
    ModuleBorderColor    string
  }
}

func NewCfg() (Cfg, error) {
  viper.SetDefault("Theme.ModuleBorderColor", "#FFFFFF")

  viper.SetConfigName("wth.yaml")
  viper.SetConfigType("yaml")
  viper.AddConfigPath("/etc/")
  viper.AddConfigPath("$XDG_CONFIG_HOME/")
  viper.AddConfigPath("$HOME/.config/")
  viper.AddConfigPath("$HOME/")
  viper.AddConfigPath(".")

  viper.SetEnvPrefix("wth")
  viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
  viper.AutomaticEnv()

  if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
      return Cfg{}, err
    }
  }

  var config Cfg
  if err := viper.Unmarshal(&config); err != nil {
    return Cfg{}, err
  }

  return config, nil
}

