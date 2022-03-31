package runner

import (
	"fmt"
)

type Command struct {
	bin  string
	args []string
}

func NewCommand(args []string) (Command, error) {

	if len(args) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}

	return Command{
		bin:  args[0],
		args: args[1:],
	}, nil
}

func (c Command) Bin() string {
	return c.bin
}

func (c Command) Args() []string {
	return c.args
}
