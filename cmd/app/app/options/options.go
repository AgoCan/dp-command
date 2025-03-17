package options

import (
	"dp-command/internal/command"
	"dp-command/internal/config"
)

const (
	_defaultConfigFile = "config/config.yaml"
)

type AppOptions struct {
	ConfFile string
	Config   *config.Config
}

func NewAppOptions() *AppOptions {
	o := &AppOptions{}
	return o
}

func (o *AppOptions) NewServer() (*command.Command, error) {
	s := command.New()
	o.loadConfig(o.ConfFile)
	s.Config = o.Config

	return s, nil
}

func (o *AppOptions) loadConfig(configFile string) {
	if configFile == "" {
		configFile = _defaultConfigFile
	}
	o.Config = config.New(configFile)
}
