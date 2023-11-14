//go:build integration
// +build integration

package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func randomTaskName(t *testing.T) string {

	t.Helper()
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var p strings.Builder
	for i := 0; i < 32; i++ {
		p.WriteByte(chars[r.Intn(len(chars))])
	}

	return p.String()
}

func TestIntegration(t *testing.T) {
	apiRoot := "http://localhost:8080"

	if os.Getenv("TODO_API_ROOT") != "" {
		apiRoot = os.Getenv("TODO_API_ROOT")
	}

	today := time.Now().Format("Jan/02")

	task := randomTaskName(t)
	taskId := ""

	t.Run("AddTask", func(t *testing.T) {
		args := []string{task}
		expOut := fmt.Sprintf("Added task %q to the list\n", task)

		var out bytes.Buffer

		if err := addAction(&out, apiRoot, args); err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		if expOut != out.String() {
			t.Errorf("expected output %q, got %q", expOut, out.String())
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		var out bytes.Buffer

		if err := listAction(&out, apiRoot); err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		outList := ""

		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				outList = scanner.Text()
				break
			}
		}

		if outList == "" {
			t.Errorf("task %q not in the list", task)
		}

		taskCompleteStatus := strings.Fields(outList)[0]

		if taskCompleteStatus != "-" {
			t.Errorf("expected status %q, got %q", "-", taskCompleteStatus)
		}

		taskId = strings.Fields(outList)[1]
	})

	vRes := t.Run("ViewTask", func(t *testing.T) {
		var out bytes.Buffer

		if err := viewAction(&out, apiRoot, taskId); err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		viewOut := strings.Split(out.String(), "\n")

		if !strings.Contains(viewOut[0], task) {
			t.Fatalf("expected task %q, got %q", task, viewOut[0])
		}

		if !strings.Contains(viewOut[1], today) {
			t.Fatalf("expected creation day/month %q, got %q", today, viewOut[1])
		}

		if !strings.Contains(viewOut[2], "No") {
			t.Fatalf("expected completed status %q, got %q", "No", viewOut[2])
		}
	})

	if !vRes {
		t.Fatalf("view task failed. Stopping integration tests")
	}

	t.Run("CompleteTask", func(t *testing.T) {
		var out bytes.Buffer

		if err := completeAction(&out, apiRoot, taskId); err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		expOut := fmt.Sprintf("item number %s marked as completed\n", taskId)

		if expOut != out.String() {
			t.Fatalf("expected output %q, got %q", expOut, out.String())
		}
	})

	t.Run("ListCompletedTask", func(t *testing.T) {
		var out bytes.Buffer

		if err := listAction(&out, apiRoot); err != nil {
			t.Errorf("expected no error, got %q", err)
		}

		outList := ""

		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				outList = scanner.Text()
				break
			}
		}

		if outList == "" {
			t.Errorf("task %q is not in the list", task)
		}

		taskCompleteStatus := strings.Fields(outList)[0]

		if taskCompleteStatus != "X" {
			t.Errorf("expected status %q, got %q", "X", taskCompleteStatus)
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		var out bytes.Buffer

		if err := delAction(&out, apiRoot, taskId); err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		expOut := fmt.Sprintf("item number %s deleted\n", taskId)

		if expOut != out.String() {
			t.Fatalf("expected output %q, got %q", expOut, out.String())
		}
	})

	t.Run("ListDeletedTask", func(t *testing.T) {
		var out bytes.Buffer
		// deleted task should not exist
		// hence try listing and check for error
		err := listAction(&out, apiRoot)

		if err == nil {
			t.Fatalf("expected %q error, got nil", ErrNotFound)
		}

		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("expected error: %q, got %q instead", ErrNotFound, err)
			}
		}

		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				t.Errorf("task %q still in list", task)
				break
			}
		}
	})
}
