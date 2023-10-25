package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shantanoo-desai/rggo/todo"
)

var todoFileName = ".todo.json"

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil

	}

	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}

func main() {
	add := flag.Bool("add", false, "Task to be included in ToDo List")
	list := flag.Bool("list", false, "List all Tasks")
	complete := flag.Int("complete", 0, "Item to be marked completed")
	delete := flag.Int("del", 0, "Item to be deleted")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Todo List Manager\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		// list current todo items if no args provided
		fmt.Print(l)

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// save updated list to file
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Add(t)
		// save updated list to file
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid Option")
		os.Exit(1)

	}
}
