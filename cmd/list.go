/*
Package cmd
Copyright © 2025 Esben Inglev <esbeninglev@gmail.com>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the current databases in the docker container",
	Long:  `Lists the current databases in the docker container`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if password == "" {
			return fmt.Errorf("SA password required. Use the -p command to set the pw")
		}

		query := "SELECT name, database_id, create_date, state_desc FROM sys.databases ORDER BY database_id"

		dockerCmd := exec.Command(
			"docker", "exec", "-i", container,
			baseCmd,
			"-S", "localhost",
			"-U", "sa",
			"-P", password,
			"-Q", query,
			"-C",
			"-h", "-1", // remove headers
			"-s", "|",
			"-W", // remove trailing whitespace
		)

		output, err := dockerCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to list databases: %w\nOutput: %s", err, string(output))
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")

		if len(lines) == 0 {
			fmt.Println("No databases found...")
			return nil
		}

		// header print
		fmt.Printf("\n%-30s %-12s %-20s %-15s\n", "DATABASE NAME", "ID", "CREATED", "STATE")
		fmt.Println(strings.Repeat("-", 80))

		for _, line := range lines {
			line = strings.TrimSpace(line)

			if line == "" {
				continue
			}

			fields := strings.Split(line, "|")
			if len(fields) >= 4 {
				name := fields[0]
				id := fields[1]
				created := fields[2]
				state := fields[3]

				// trim name if it's too long
				if len(name) > 30 {
					name = name[:27] + "..."
				}

				// only get date part
				if len(created) > 10 {
					created = created[:10]
				}

				fmt.Printf("%-30s %-12s %-20s %-15s\n", name, id, created, state)
			}
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
