//go:build inmemory
// +build inmemory

package pomodoro_test

import (
	"testing"

	"github.com/shantanoo-desai/rggo/pomo/pomodoro"
	"github.com/shantanoo-desai/rggo/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
