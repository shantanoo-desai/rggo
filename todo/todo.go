package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// ToDo Item
type item struct {
	Done        bool
	Task        string
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List of ToDos
type List []item

// Add item to a List of ToDos
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Mark an item as complete within a List of ToDos
func (l *List) Complete(i int) error {
	ls := *l

	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	// Adjust 0-base index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete a ToDo from the List
func (l *List) Delete(i int) error {
	ls := *l

	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save current List to a JSON file locally
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)

	if err != nil {
		return fmt.Errorf("Marshalling JSON error: %v", err)
	}

	return os.WriteFile(filename, js, 0644)
}

// Read a JSON list of Todos to memory
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)

}

// Formatted output by satisfying the Stringer Interface in "fmt"
func (l *List) String() string {
	formatted := ""

	for k, t := range *l {
		prefix := " "
		if t.Done {
			prefix = "X "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}
	return formatted
}
