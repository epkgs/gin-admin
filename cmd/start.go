package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gin-admin/internal/app"
	"gin-admin/internal/configs"

	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {

	// startCmd represents the start command
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile, _ := cmd.Flags().GetString("config")

			if daemon, _ := cmd.Flags().GetBool("daemon"); daemon {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					fmt.Printf("failed to get absolute path for command: %s \n", err.Error())
					return err
				}

				args := []string{"start"}
				args = append(args, "-c", configFile)
				fmt.Printf("execute command: %s %s \n", bin, strings.Join(args, " "))
				command := exec.Command(bin, args...)

				// Redirect stdout and stderr to log file
				stdLogFile := fmt.Sprintf("%s.log", cmd.Root().Name())
				file, err := os.OpenFile(stdLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					fmt.Printf("failed to open log file: %s \n", err.Error())
					return err
				}
				defer file.Close()

				command.Stdout = file
				command.Stderr = file

				err = command.Start()
				if err != nil {
					fmt.Printf("failed to start daemon thread: %s \n", err.Error())
					return err
				}

				// Don't wait for the command to finish
				// The main process will exit, allowing the daemon to run independently
				fmt.Printf("Service %s daemon thread started successfully\n", configs.C.AppName)

				pid := command.Process.Pid
				_ = os.WriteFile(fmt.Sprintf("%s.lock", cmd.Root().Name()), []byte(fmt.Sprintf("%d", pid)), 0666)
				fmt.Printf("service %s daemon thread started with pid %d \n", configs.C.AppName, pid)
				os.Exit(0)
			}

			err := app.Run(context.Background(), configFile)
			if err != nil {
				panic(err)
			}
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "config.yaml", "Config file")
	cmd.Flags().BoolP("daemon", "d", false, "Run as a daemon")

	return cmd
}
