/*
Copyright © 2023 Shantanoo 'Shan' Desai

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View details about a single item",
	SilenceUsage: true,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")

		return viewAction(os.Stdout, apiRoot, args[0])
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func viewAction(out io.Writer, apiRoot, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: item id must be a number", ErrNotNumber)
	}

	i, err := getOne(apiRoot, id)
	if err != nil {
		return err
	}

	return printOne(out, i)
}

func printOne(out io.Writer, i item) error {
	w := tabwriter.NewWriter(out, 14, 2, 0, ' ', 0)

	fmt.Fprintf(w, "Task:\t%s\n", i.Task)
	fmt.Fprintf(w, "Created at:\t%s\n", i.CreatedAt.Format(timeFormat))

	if i.Done {
		fmt.Fprintf(w, "Completed:\t%s\n", "Yes")
		fmt.Fprintf(w, "CompletedAt:\t%s\n", i.CompletedAt.Format(timeFormat))
		return w.Flush()
	}

	fmt.Fprintf(w, "Completed:\t%s\n", "No")
	return w.Flush()
}
