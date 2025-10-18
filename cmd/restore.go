/*
Package cmd
Copyright © 2025 Esben Inglev <esbeninglev@gmail.com>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores a database form a .bak file",
	Long:  `Restores a database from a .bak file`,
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

		if database == "" {
			return fmt.Errorf("database name must be provoded. Use the -d flag to set the database name")
		}

		query := fmt.Sprintf("RESTORE DATABASE [%s] FROM DISK = N'/var/opt/mssql/backup/%s' WITH REPLACE", database, file)

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
			return fmt.Errorf("failed to restore database: %w\nOutput: %s", err, string(output))
		}

		fmt.Printf("✓ Database '%s' restored successfully from %s\n", database, file)
		fmt.Print(string(output))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
