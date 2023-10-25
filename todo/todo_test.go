package todo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/shantanoo-desai/rggo/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"

	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, Got %q instead\n", taskName, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"

	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, Got %q instead\n", taskName, l[0].Task)
	}

	l.Complete(1)

	if !l[0].Done {
		fmt.Errorf("New Task should be Marked Completed\n")
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}

	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}

	for _, v := range tasks {
		l.Add(v)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("Expected %q, Got %q instead\n", tasks[0], l[0].Task)
	}

	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list length %d, Got %d instead\n", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, Got %q instead\n", tasks[2], l[1].Task)
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New Task"
	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Expected %q, Got %q instead\n", taskName, l1[0].Task)
	}

	tempFile, err := os.CreateTemp("", "")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if err := l1.Save(tempFile.Name()); err != nil {
		t.Fatalf("Error Saving List to File: %s", err)
	}

	if err := l2.Get(tempFile.Name()); err != nil {
		t.Fatalf("Error Getting List from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match Task %q", l1[0].Task, l2[0].Task)
	}
}
