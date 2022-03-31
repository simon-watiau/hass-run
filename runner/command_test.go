package runner

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
}

func (suite *CommandTestSuite) TestValidCommand() {
	command, err := NewCommand([]string{"bash", "-c", "ls && ls"})
	suite.Nil(err)
	suite.Equal("bash", command.Bin())
	suite.Equal([]string{"-c", "ls && ls"}, command.Args())
}

func (suite *CommandTestSuite) TestWithoutArguments() {
	command, err := NewCommand([]string{"ls"})
	suite.Nil(err)
	suite.Equal("ls", command.Bin())
	suite.Equal([]string{}, command.Args())
}

func (suite *CommandTestSuite) TestEmptyCommand() {
	_, err := NewCommand([]string{})
	suite.NotNil(err)

}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
