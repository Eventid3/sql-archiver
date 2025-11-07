/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/Eventid3/sql-archiver/interactive"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

func main() {
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
		container := viper.GetString("container")
		user := viper.GetString("user")
		password := viper.GetString("password")

		p := tea.NewProgram(interactive.InitialModelWithConfig(container, user, password), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}
