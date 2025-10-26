/*
Package cmd
Copyright © 2025 Esben Inglev <esbeninglev@gmail.com>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/Eventid3/sql-archiver/mssql"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the current databases in the docker container",
	Long:  `Lists the current databases in the docker container`,

	RunE: func(cmd *cobra.Command, args []string) error {
		dbList, err := mssql.GetDatabases(container, "sa", password)
		if err != nil {
			return fmt.Errorf("error retrieving databases: %v", err)
		}

		// header print
		fmt.Printf("\n%-30s %-12s %-20s %-15s\n", "DATABASE NAME", "ID", "CREATED", "STATE")
		fmt.Println(strings.Repeat("-", 80))

		for _, item := range dbList {
			if len(item.Name) > 28 {
				item.Name = item.Name[:25] + "..."
			}

			fmt.Printf("%-30s %-12s %-20s %-15s\n", item.Name, item.ID, item.Created, item.State)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
