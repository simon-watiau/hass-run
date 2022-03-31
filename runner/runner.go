package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const CommandFailedExitCode = -10

type Hass interface {
	UpdateState(json string) error
}

type Runner struct {
	command   Command
	hass      Hass
	output    string
	running   bool
	exitCode  int
	startedAt time.Time
	updatedAt time.Time
	endedAt   time.Time
	duration  time.Duration
}

type CommandRun interface {
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	Kill() error
}

type commandRun struct {
	*exec.Cmd
}

func (c *commandRun) Kill() error {
	if c.Process == nil {
		return errors.New("failed to kill a non running process")
	}

	return c.Process.Kill()
}

var Executor = func(cmd string, args []string) CommandRun {
	execCmd := exec.Command(cmd, args...)
	return &commandRun{execCmd}
}

func NewRunner(command Command, hass Hass) *Runner {
	return &Runner{
		command: command,
		hass:    hass,
	}
}

func (r *Runner) Run() {
	r.output = ""
	r.running = true
	r.startedAt = time.Now()
	r.endedAt = time.Time{}
	r.updatedAt = time.Time{}
	r.duration = 0

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM)

	r.Notify()
	defer r.Notify()

	cmd := Executor(r.command.Bin(), r.command.Args())

	go func() {
		select {
		case <-cancelChan:
			cmd.Kill()
		case <-context.Done():
		}
	}()

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Printf("Failed to acquire STDOUT pipe: %s", err.Error())
		r.appendOutput(err.Error() + "\n")
		r.Notify()
		return
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		log.Printf("Failed to acquire STDERR pipe: %s", err.Error())
		r.appendOutput(err.Error() + "\n")
		return
	}

	var wg sync.WaitGroup

	go r.ReadStream(&wg, &r.output, stdout)

	go r.ReadStream(&wg, &r.output, stderr)

	err = cmd.Start()

	if err != nil {
		log.Printf("Failed to run command: %s", err.Error())
		r.appendOutput(err.Error() + "\n")
		return
	}

	wg.Wait()

	err = cmd.Wait()
	r.running = false
	r.endedAt = time.Now()

	if exitCode, ok := err.(*exec.ExitError); ok {
		log.Printf(
			"Command failed with status code: %d",
			exitCode.ExitCode(),
		)

		r.exitCode = exitCode.ExitCode()

		return
	}

	if err != nil {
		log.Printf(
			"Command failed: %s",
			err.Error(),
		)

		r.exitCode = CommandFailedExitCode
		r.appendOutput(err.Error() + "\n")
	}
}

func (r *Runner) ReadStream(
	wg *sync.WaitGroup,
	output *string,
	reader io.ReadCloser,
) {
	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			log.Println(scanner.Text())
			r.appendOutput(scanner.Text() + "\n")
			r.Notify()
		}
		wg.Done()
	}()
}

func (r *Runner) appendOutput(content string) {
	r.output += content
	r.updatedAt = time.Now()
}

func (r *Runner) Notify() {
	var state string
	if r.running {
		state = "running"
	} else {
		if r.exitCode == 0 {
			state = "success"
		} else {
			state = "failure"
		}
	}

	payload := Payload{
		State: state,
		Attributes: Attributes{
			Output:    r.output,
			ExitCode:  r.exitCode,
			StartedAt: r.startedAt,
			UpdatedAt: r.updatedAt,
			EndedAt:   r.endedAt,
		},
	}

	if (r.endedAt != time.Time{}) {
		payload.Attributes.Duration = int(r.endedAt.Sub(r.startedAt).Seconds())
	}

	bytes, err := json.Marshal(payload)

	if err != nil {
		log.Printf(
			"Failed to marshal payload: %s",
			err.Error(),
		)
		return
	}

	err = r.hass.UpdateState(string(bytes))

	if err != nil {
		log.Printf(
			"Failed to publish update: %s",
			err.Error(),
		)
	}
}
