package cmd

import (
	"fmt"
	"log"

	"github.com/sevlyar/go-daemon"
	"github.com/simon-watiau/mqtt-run/hass"
	"github.com/simon-watiau/mqtt-run/pid"
	"github.com/simon-watiau/mqtt-run/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:        "run [flags] [entity] [PIDFile] -- [command]",
	Short:      "Run a command",
	Args:       cobra.MinimumNArgs(3),
	ArgAliases: []string{"entity", "PIDFile", "command"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validate(cmd, args)
		if err != nil {
			return err
		}
		err = run(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("host", "f", "", "HomeAssistant host (e.g https://hass.fr)")
	runCmd.Flags().StringP("bearer", "b", "", "Bearer token for HomeAssistant")
	runCmd.Flags().BoolP("nodaemon", "n", false, "Disable daemon for debug")

	viper.BindPFlags(runCmd.Flags())

	viper.SetConfigName("mqtt-run")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("/etc")
}

func validate(cmd *cobra.Command, args []string) error {
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	err = hass.ValidateEntityName(args[0])

	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	err = hass.ValidateHostAndBearer(
		viper.GetString("host"),
		viper.GetString("bearer"),
	)

	if err != nil {
		return fmt.Errorf("invalid host/bearer: %w", err)
	}

	err = pid.ValidatePIDFile(
		args[1],
	)

	if err != nil {
		return fmt.Errorf("invalid PID file: %w", err)
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	context := &daemon.Context{
		PidFileName: args[1],
		PidFilePerm: 0644,
	}

	if !viper.GetBool("nodaemon") {
		child, err := context.Reborn()

		if err != nil {
			return fmt.Errorf("failed to spawn daemon: %w", err)
		}

		if child != nil {
			log.Printf("Child process started")
			// exit parent
			return nil
		}

		defer context.Release()
	}

	command, err := runner.NewCommand(args[2:])

	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}

	hass := hass.NewHass(
		viper.GetString("bearer"),
		viper.GetString("host"),
		args[0],
	)

	cmdRunner := runner.NewRunner(
		command,
		hass,
	)

	cmdRunner.Run()

	return nil
}
