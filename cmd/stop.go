package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func StopCmd() *cobra.Command {

	// stopCmd represents the stop command
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop server",
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := cmd.Root().Name()
			lockFile := fmt.Sprintf("%s.lock", appName)
			pid, err := os.ReadFile(lockFile)
			if err != nil {
				return err
			}

			command := exec.Command("kill", string(pid))
			err = command.Start()
			if err != nil {
				return err
			}

			err = os.Remove(lockFile)
			if err != nil {
				return fmt.Errorf("can't remove %s.lock. %s", appName, err.Error())
			}

			fmt.Printf("service %s stopped \n", appName)
			return nil
		},
	}

	return cmd
}
