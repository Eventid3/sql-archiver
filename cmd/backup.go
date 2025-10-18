/*
Package cmd
Copyright © 2025 Esben Inglev <esbeninglev@gmail.com>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var timestamp bool

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Makes a backup of the specified database",
	Long:  `Makes a backup of the specified database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if password == "" {
			return fmt.Errorf("SA password required. Use the -p command to set the pw")
		}

		if database == "" {
			return fmt.Errorf("database name must be provided. Use the -d flag to set the database name")
		}

		if file == "" {
			return fmt.Errorf("filename must be provided. Use the -f command to set the filename")
		}

		if len(file) < 5 || file[len(file)-4:] != ".bak" {
			return fmt.Errorf("filename must end with .bak")
		}

		t := time.Now()
		timeString := t.Format("2006-01-02T15-04-05")

		if timestamp {
			file = file[:len(file)-4] + "_" + timeString + ".bak"
		}

		query := fmt.Sprintf("BACKUP DATABASE [%s] TO DISK = N'/var/opt/mssql/backup/%s' WITH STATS = 10", database, file)

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
			return fmt.Errorf("failed to backup database: %w\nOutput: %s", err, string(output))
		}

		fmt.Printf("✓ Database '%s' backed up successfully to %s\n", database, file)
		fmt.Print(string(output))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// local flags
	backupCmd.Flags().BoolVarP(&timestamp, "timestamp", "t", false, "Add a timestamp to the .bak file")
}
