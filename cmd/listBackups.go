/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// listBackupsCmd represents the listBackups command
var listBackupsCmd = &cobra.Command{
	Use:   "listBackups",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if password == "" {
			return fmt.Errorf("SA password required. Use the -p command to set the pw")
		}

		dockerCmd := exec.Command(
			"docker", "exec", "-i", container,
			"ls", "-lhgG", "/var/opt/mssql/backup",
		)

		output, err := dockerCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to list backup files: %w\nOutput: %s", err, string(output))
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")

		if len(lines) == 0 {
			fmt.Println("No backup files found...")
			return nil
		}

		// header print
		fmt.Printf("\n%-10s %-20s %s\n", "SIZE", "DATE", "FILE NAME")
		fmt.Println(strings.Repeat("-", 80))

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "total") {
				continue
			}

			// Use Fields instead of Split to handle multiple spaces
			fields := strings.Fields(line)

			// ls -lhgG output format:
			// [0] permissions, [1] links, [2] size, [3] month, [4] day, [5] time/year, [6] filename
			if len(fields) >= 7 {
				size := fields[2]
				date := fmt.Sprintf("%s %s %s", fields[3], fields[4], fields[5])
				name := fields[6]

				// Handle filenames with spaces (join remaining fields)
				if len(fields) > 7 {
					name = strings.Join(fields[6:], " ")
				}

				fmt.Printf("%-10s %-20s %s\n", size, date, name)
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listBackupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listBackupsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listBackupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
