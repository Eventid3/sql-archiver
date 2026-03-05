/*
Package cmd
Copyright © 2025 Esben Inglev <esbeninglev@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Eventid3/sql-archiver/interactive"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ignoreConfig bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sql-archiver",
	Short: "Tool to restore and backup MSSQL databases",
	Long: `Tool to restore and backup MSSQL databases.
	Restores and backups are done with .bak files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ignoreConfig {
			p := tea.NewProgram(interactive.InitialModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		}

		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.config/sql-archiver/")

		err := viper.ReadInConfig()

		if err != nil {
			p := tea.NewProgram(interactive.InitialModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		} else {
			host := viper.GetString("host")
			user := viper.GetString("user")
			password := viper.GetString("password")

			p := tea.NewProgram(interactive.InitialModelWithConfig(host, user, password), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		}
	},
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
	rootCmd.PersistentFlags().BoolVarP(&ignoreConfig, "manual", "m", false, "Manual login page")
}
