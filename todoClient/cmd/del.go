/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "Deletes an item from the list",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")

		return delAction(os.Stdout, apiRoot, args[0])
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// delCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// delCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func delAction(out io.Writer, apiRoot, arg string) error {
	id, err := strconv.Atoi(arg)

	if err != nil {
		return fmt.Errorf("%w: item id must be a number", ErrNotNumber)
	}

	if err := deleteItem(apiRoot, id); err != nil {
		return err
	}

	return printDel(out, id)
}

func printDel(out io.Writer, id int) error {
	_, err := fmt.Fprintf(out, "item number %d deleted\n", id)
	return err
}
