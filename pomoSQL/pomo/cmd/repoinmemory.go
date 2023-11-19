//go:build inmemory
// +build inmemory

package cmd

import (
	"github.com/shantanoo-desai/rggo/pomo/pomodoro"
	"github.com/shantanoo-desai/rggo/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
