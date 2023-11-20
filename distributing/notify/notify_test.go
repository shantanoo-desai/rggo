//go:build !integration
// +build !integration

package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		s Severity
	}{
		{SeverityLow},
		{SeverityNormal},
		{SeverityUrgent},
	}

	for _, tc := range testCases {
		name := tc.s.String()
		expMessage := "Message"
		expTitle := "Title"
		t.Run(name, func(t *testing.T) {
			n := New(expTitle, expMessage, tc.s)
			if n.message != expMessage {
				t.Errorf("expected %q, got %q instead\n", expMessage, n.message)
			}

			if n.title != expTitle {
				t.Errorf("expected %q, got %q instead\n", expTitle, n.title)
			}

			if n.severity != tc.s {
				t.Errorf("expected %q, got %q instead\n", tc.s, n.severity)
			}
		})
	}
}

func TestSeverityString(t *testing.T) {
	testCases := []struct {
		s   Severity
		exp string
		os  string
	}{
		{SeverityLow, "low", "linux"},
		{SeverityNormal, "normal", "linux"},
		{SeverityUrgent, "critical", "linux"},
		{SeverityLow, "Low", "darwin"},
		{SeverityNormal, "Normal", "darwin"},
		{SeverityUrgent, "Critical", "darwin"},
		{SeverityLow, "Info", "windows"},
		{SeverityNormal, "Warning", "windows"},
		{SeverityUrgent, "Error", "windows"},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%s%d", tc.os, tc.s)
		t.Run(name, func(t *testing.T) {
			if runtime.GOOS != tc.os {
				t.Skip("Skipped: not OS", runtime.GOOS)
			}
			sev := tc.s.String()

			if sev != tc.exp {
				t.Errorf("expected %q, got %q instead\n", tc.exp, sev)
			}
		})
	}
}

func mockCmd(exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	cmdName := ""

	switch runtime.GOOS {
	case "linux":
		cmdName = "notify-send"
	case "darwin":
		cmdName = "terminal-notifier"
	case "windows":
		cmdName = "powershell"
	}

	if strings.Contains(os.Args[2], cmdName) {
		os.Exit(0)
	}

	os.Exit(1)
}

func TestSend(t *testing.T) {
	n := New("test title", "test msg", SeverityNormal)

	command = mockCmd

	err := n.Send()

	if err != nil {
		t.Error(err)
	}
}
