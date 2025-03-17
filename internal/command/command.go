package command

import (
	"fmt"

	"dp-command/internal/config"
)

type Command struct {
	Config *config.Config
}

func New() *Command {
	return &Command{}
}

func (s *Command) Run() {
	fmt.Println("run command")
}
