package pomodoro_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/shantanoo-desai/rggo/pomo/pomodoro"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name   string
		input  [3]time.Duration
		expect pomodoro.IntervalConfig
	}{
		{
			name: "Default",
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   25 * time.Minute,
				ShortBreakDuration: 5 * time.Minute,
				LongBreakDuration:  15 * time.Minute,
			},
		},
		{
			name: "SingleInput",
			input: [3]time.Duration{
				20 * time.Minute,
			},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   20 * time.Minute,
				ShortBreakDuration: 5 * time.Minute,
				LongBreakDuration:  15 * time.Minute,
			},
		},
		{
			name: "MultipleInput",
			input: [3]time.Duration{
				20 * time.Minute,
				10 * time.Minute,
				12 * time.Minute,
			},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   20 * time.Minute,
				ShortBreakDuration: 10 * time.Minute,
				LongBreakDuration:  12 * time.Minute,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var repo pomodoro.Repository
			config := pomodoro.NewConfig(
				repo,
				tc.input[0],
				tc.input[1],
				tc.input[2],
			)

			if config.PomodoroDuration != tc.expect.PomodoroDuration {
				t.Errorf("expected PomodoroDuration: %q, got %q\n", tc.expect.PomodoroDuration, config.PomodoroDuration)
			}

			if config.ShortBreakDuration != tc.expect.ShortBreakDuration {
				t.Errorf("expected ShortBreakDuration: %q, got %q\n", tc.expect.ShortBreakDuration, config.ShortBreakDuration)
			}

			if config.LongBreakDuration != tc.expect.LongBreakDuration {
				t.Errorf("expected LongBreakDuration: %q, got %q\n", tc.expect.LongBreakDuration, config.LongBreakDuration)
			}
		})
	}
}

func TestGetInterval(t *testing.T) {
	repo, cleanup := getRepo(t)

	defer cleanup()

	const duration = 1 * time.Millisecond
	config := pomodoro.NewConfig(repo, 3*duration, duration, 2*duration)

	for i := 1; i <= 16; i++ {
		var (
			expCategory string
			expDuration time.Duration
		)

		switch {
		case i%2 != 0:
			expCategory = pomodoro.CategoryPomodoro
			expDuration = 3 * duration
		case i%8 == 0:
			expCategory = pomodoro.CategoryLongBreak
			expDuration = 2 * duration
		case i%2 == 0:
			expCategory = pomodoro.CategoryShortBreak
			expDuration = duration
		}

		testName := fmt.Sprintf("%s%d", expCategory, i)
		t.Run(testName, func(t *testing.T) {
			res, err := pomodoro.GetInterval(config)

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			noop := func(pomodoro.Interval) {}

			if err := res.Start(context.Background(), config, noop, noop, noop); err != nil {
				t.Fatal(err)
			}

			if res.Category != expCategory {
				t.Errorf("expected category %q, got %q", expCategory, res.Category)
			}

			if res.PlannedDuration != expDuration {
				t.Errorf("expected PlannedDuration %q, got %q", expDuration, res.PlannedDuration)
			}

			if res.State != pomodoro.StateNotStarted {
				t.Errorf("expected State=%q, got %q", pomodoro.StateNotStarted, res.State)
			}

			ui, err := repo.ByID(res.ID)
			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			if ui.State != pomodoro.StateDone {
				t.Errorf("expected State= %q, got %q", pomodoro.StateDone, res.State)
			}
		})
	}
}

func TestPause(t *testing.T) {
	const duration = 2 * time.Second

	repo, cleanup := getRepo(t)
	defer cleanup()

	config := pomodoro.NewConfig(repo, duration, duration, duration)

	testCases := []struct {
		name        string
		start       bool
		expState    int
		expDuration time.Duration
	}{

		{
			name:        "NotStarted",
			start:       false,
			expState:    pomodoro.StateNotStarted,
			expDuration: 0,
		},
		{
			name:        "Paused",
			start:       true,
			expState:    pomodoro.StatePaused,
			expDuration: duration / 2,
		},
	}

	expError := pomodoro.ErrIntervalNotRunning

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			i, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}

			start := func(pomodoro.Interval) {}
			end := func(pomodoro.Interval) {
				t.Errorf("end callback should not be executed")
			}

			periodic := func(i pomodoro.Interval) {
				if err := i.Pause(config); err != nil {
					t.Fatal(err)
				}
			}

			if tc.start {
				if err := i.Start(ctx, config, start, periodic, end); err != nil {
					t.Fatal(err)
				}
			}

			i, err = pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}

			err = i.Pause(config)
			if err != nil {
				if !errors.Is(err, expError) {
					t.Fatalf("expected error %q, got %q", expError, err)
				}
			}

			if err == nil {
				t.Errorf("expected error %q, got nil", expError)
			}

			i, err = repo.ByID(i.ID)

			if err != nil {
				t.Fatal(err)
			}

			if i.State != tc.expState {
				t.Errorf("expected state %d, got %d", tc.expState, i.State)
			}

			if i.ActualDuration != tc.expDuration {
				t.Errorf("expected duration %q, got %q", tc.expDuration, i.ActualDuration)
			}

			cancel()
		})
	}
}

func TestStart(t *testing.T) {
	const duration = 2 * time.Second

	repo, cleanup := getRepo(t)
	defer cleanup()

	config := pomodoro.NewConfig(repo, duration, duration, duration)

	testCases := []struct {
		name        string
		cancel      bool
		expState    int
		expDuration time.Duration
	}{
		{
			name:        "Finish",
			cancel:      false,
			expState:    pomodoro.StateDone,
			expDuration: duration,
		},
		{
			name:        "Cancel",
			cancel:      true,
			expState:    pomodoro.StateCancelled,
			expDuration: duration / 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			i, err := pomodoro.GetInterval(config)

			if err != nil {
				t.Fatal(err)
			}

			start := func(i pomodoro.Interval) {
				if i.State != pomodoro.StateRunning {
					t.Errorf("expected state %d, got %d", pomodoro.StateRunning, i.State)
				}

				if i.ActualDuration >= i.PlannedDuration {
					t.Errorf("expected ActualDuration %q, less than Planned %q", i.ActualDuration, i.PlannedDuration)
				}
			}

			end := func(i pomodoro.Interval) {
				if i.State != tc.expState {
					t.Errorf("expected state %d, got %d", tc.expState, i.State)
				}

				if tc.cancel {
					t.Errorf("end callback should not be executed")
				}
			}

			periodic := func(i pomodoro.Interval) {
				if i.State != pomodoro.StateRunning {
					t.Errorf("expected state %d, got %d", pomodoro.StateRunning, i.State)
				}

				if tc.cancel {
					cancel()
				}
			}

			if err := i.Start(ctx, config, start, periodic, end); err != nil {
				t.Fatal(err)
			}

			i, err = repo.ByID(i.ID)
			if err != nil {
				t.Fatal(err)
			}

			if i.State != tc.expState {
				t.Errorf("expected state %d, got %d", tc.expState, i.State)
			}

			if i.ActualDuration != tc.expDuration {
				t.Errorf("expected ActualDuration %q, got %q", tc.expDuration, i.ActualDuration)
			}

			cancel()
		})

	}
}
