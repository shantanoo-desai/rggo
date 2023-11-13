/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate bash completions for your command",
	Long: `to load completions run:
	source < (pScan completion)

	to load completions automatically add the following line to ~/.bashrc file:
	source < (pScan completion)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completeAction(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func completeAction(out io.Writer) error {
	return rootCmd.GenBashCompletion(out)
}
