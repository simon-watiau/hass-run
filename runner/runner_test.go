package runner

import (
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HassMock struct {
	mock.Mock
}

func (h *HassMock) UpdateState(json string) error {
	return h.Called(json).Error(0)
}

var CommandBin = "cmd"
var CommandArgs = []string{"a", "b c"}

type CmdMock struct {
	mock.Mock
}

func (c *CmdMock) StderrPipe() (io.ReadCloser, error) {
	args := c.Called()
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (c *CmdMock) StdoutPipe() (io.ReadCloser, error) {
	args := c.Called()
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (c *CmdMock) Start() error {
	return c.Called().Error(0)
}

func (c *CmdMock) Wait() error {
	return c.Called().Error(0)
}

func (c *CmdMock) Kill() error {
	return c.Called().Error(0)
}

type RunnerTestSuite struct {
	suite.Suite
	hassMock *HassMock
	cmdMock  *CmdMock
	runner   *Runner
}

func (suite *RunnerTestSuite) SetupTest() {
	suite.hassMock = &HassMock{}
	suite.cmdMock = &CmdMock{}

	Executor = func(cmd string, args []string) CommandRun {
		suite.Equal(CommandBin, cmd)
		suite.Equal(CommandArgs, args)
		return suite.cmdMock
	}

	command, err := NewCommand(append([]string{CommandBin}, CommandArgs...))
	suite.Nil(err)

	suite.runner = NewRunner(
		command,
		suite.hassMock,
	)
}

func (suite *RunnerTestSuite) TestValidCommand() {
	startDate := time.Now()
	output1Date := time.Now().Add(5 * time.Second)
	output2Date := time.Now().Add(10 * time.Second)
	endDate := time.Now().Add(20 * time.Second)

	notified := make(chan int)

	monkey.Patch(time.Now, func() time.Time {
		return startDate
	})

	stdoutReader, stdoutWriter := io.Pipe()

	suite.cmdMock.On("StdoutPipe").Return(stdoutReader, nil)

	stderrReader, stderrWriter := io.Pipe()
	suite.cmdMock.On("StderrPipe").Return(stderrReader, nil)

	suite.cmdMock.On("Start").Return(nil)

	waitChan := make(chan time.Time)
	suite.cmdMock.On("Wait").Return(errors.New("FAILED")).WaitFor = waitChan

	go func() {
		<-notified

		monkey.Patch(time.Now, func() time.Time {
			return output1Date
		})
		stdoutWriter.Write([]byte("Hello world\n"))

		<-notified

		monkey.Patch(time.Now, func() time.Time {
			return output2Date
		})
		stderrWriter.Write([]byte("An error\n"))

		<-notified

		monkey.Patch(time.Now, func() time.Time {
			return endDate
		})

		stdoutReader.Close()
		stderrReader.Close()
		close(waitChan)
	}()

	suite.hassMock.On(
		"UpdateState",
		suite.Payload(Payload{
			State: "running",
			Attributes: Attributes{
				Output:    "",
				StartedAt: startDate,
			},
		}),
	).Return(nil).Once().Run(func(args mock.Arguments) { notified <- 1 })

	suite.hassMock.On(
		"UpdateState",
		suite.Payload(Payload{
			State: "running",
			Attributes: Attributes{
				Output:    "Hello world\n",
				StartedAt: startDate,
				UpdatedAt: output1Date,
			},
		}),
	).Return(nil).Once().Run(func(args mock.Arguments) { notified <- 1 })

	suite.hassMock.On(
		"UpdateState",
		suite.Payload(Payload{
			State: "running",
			Attributes: Attributes{
				Output:    "Hello world\nAn error\n",
				StartedAt: startDate,
				UpdatedAt: output2Date,
			},
		}),
	).Return(nil).Once().Run(func(args mock.Arguments) { notified <- 1 })

	suite.hassMock.On(
		"UpdateState",
		suite.Payload(Payload{
			State: "failure",
			Attributes: Attributes{
				Output:    "Hello world\nAn error\nFAILED\n",
				ExitCode:  -10,
				StartedAt: startDate,
				UpdatedAt: endDate,
				EndedAt:   endDate,
				Duration:  20,
			},
		}),
	).Return(nil).Once()

	suite.runner.Run()
}

func (suite *RunnerTestSuite) Payload(payload Payload) string {
	bytes, err := json.Marshal(payload)
	suite.Nil(err)
	return string(bytes)
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerTestSuite))
}
