/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"github.com/simon-watiau/hass-run/hass"
	"github.com/simon-watiau/hass-run/pid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var killCmd = &cobra.Command{
	Use:        "kill [entity] [PIDFile]",
	Short:      "Kill a running command",
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"entity", "PIDFile"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateKillConfig(cmd, args)
		if err != nil {
			return err
		}
		err = kill(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(killCmd)

	killCmd.Flags().StringP("host", "f", "", "HomeAssistant host (e.g https://hass.fr)")
	killCmd.Flags().StringP("bearer", "b", "", "Bearer token for HomeAssistant")

	viper.BindPFlags(killCmd.Flags())

	viper.SetConfigName("hass-run")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("/etc")
}

func validateKillConfig(cmd *cobra.Command, args []string) error {
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

func kill(cmd *cobra.Command, args []string) error {
	context := &daemon.Context{
		PidFileName: args[1],
		PidFilePerm: 0644,
		Args:        os.Args,
	}

	child, err := context.Search()

	if err != nil {
		return fmt.Errorf("failed to look for running command: %w", err)
	}

	err = child.Signal(syscall.SIGTERM)

	if err != nil {
		return fmt.Errorf("failed to kill running command: %w", err)
	}

	return nil
}
