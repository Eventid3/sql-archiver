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

var verbose bool

// inspectFileCmd represents the inspectFile command
var inspectFileCmd = &cobra.Command{
	Use:   "inspectFile",
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

		if file == "" {
			return fmt.Errorf("filename must be provoded. Use the -f command to set the filename")
		}

		if len(file) < 5 || file[len(file)-4:] != ".bak" {
			return fmt.Errorf("filename must be of type .bak")
		}

		query := fmt.Sprintf("RESTORE FILELISTONLY FROM DISK = N'/var/opt/mssql/backup/%s'", file)

		dockerCmd := exec.Command(
			"docker", "exec", "-i", container,
			baseCmd,
			"-S", "localhost",
			"-U", "sa",
			"-P", password,
			"-Q", query,
			"-C",
		)

		output, err := dockerCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error listing entries from .bak file: %w\nOutput: %s", err, output)
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")

		fmt.Printf("\n%-30s %-12s %-12s\n", "DATABASE NAME", "SIZE", "BACKUPSIZE")
		fmt.Println(strings.Repeat("-", 80))

		columns := make(map[int]string)

		logicalNameCol := 0
		typeCol := 2
		sizeCol := 4
		backupSizeCol := 12

		for i, line := range lines {
			if i == 0 {
				fields := strings.Fields(line)
				for j, field := range fields {
					columns[j] = field
				}
				continue
			}

			if i == 1 {
				continue
			}

			fields := strings.Fields(line)
			if !verbose && len(fields) >= typeCol && fields[typeCol] == "D" {
				fmt.Printf("%-30s %-12s %-12s\n", fields[logicalNameCol], fields[sizeCol], fields[backupSizeCol])
			} else {
				fmt.Println(strings.Repeat("-", 80))
				if verbose {
					for j, field := range fields {
						fmt.Printf("%v: %s\n", columns[j], field)
					}
				}
				fmt.Println(strings.Repeat("-", 80))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	inspectFileCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "For detailed print")
}
