/*
Package cmd
Copyright © 2025 Esben Inglev <esbeninglev@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	container string
	password  string
	baseCmd   string = "/opt/mssql-tools18/bin/sqlcmd"

	database string
	file     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sql-archiver",
	Short: "Tool to restore and backup MSSQL databases running in docker",
	Long: `Tool to restore and backup MSSQL databases running in docker.
	Restores and backups are done with .bak files.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// global flags
	rootCmd.PersistentFlags().StringVarP(&container, "container", "c", "mssql", "Docker container name")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for SA user")
	rootCmd.PersistentFlags().StringVarP(&database, "database", "d", "", "Database name for restore or backup")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Filename to restore form or backup to")
}
