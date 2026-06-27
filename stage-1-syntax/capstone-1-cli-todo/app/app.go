// Package app parses CLI commands and connects todo domain logic with JSON storage.
//
// import 路径：just-go/stage-1-syntax/capstone-1-cli-todo/app
package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"just-go/stage-1-syntax/capstone-1-cli-todo/store"
	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

var ErrInvalidCommand = errors.New("invalid todo command")

const dataFileEnv = "JUST_GO_TODO_FILE"

// DefaultPath returns the JSON data file used by the CLI.
func DefaultPath() string {
	if path := os.Getenv(dataFileEnv); path != "" {
		return path
	}
	return ".just-go-todos.json"
}

// Run executes one CLI command.
func Run(args []string, stdout, stderr io.Writer) error {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(stdout, Usage())
		return nil
	}

	switch args[0] {
	case "add", "list", "done", "delete", "clear":
		// Valid commands continue below and may need storage.
	default:
		fmt.Fprint(stderr, Usage())
		return fmt.Errorf("%w: %s", ErrInvalidCommand, args[0])
	}

	st := store.Store{Path: DefaultPath()}
	list, err := st.Load()
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	switch args[0] {
	case "add":
		title := strings.Join(args[1:], " ")
		task, err := list.Add(title, now)
		if err != nil {
			return err
		}
		if err := st.Save(list); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "added #%d %s\n", task.ID, task.Title)
	case "list":
		fmt.Fprint(stdout, Render(list))
	case "done":
		id, err := parseID(args)
		if err != nil {
			return err
		}
		task, err := list.MarkDone(id, now)
		if err != nil {
			return err
		}
		if err := st.Save(list); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "done #%d %s\n", task.ID, task.Title)
	case "delete":
		id, err := parseID(args)
		if err != nil {
			return err
		}
		task, err := list.Delete(id)
		if err != nil {
			return err
		}
		if err := st.Save(list); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "deleted #%d %s\n", task.ID, task.Title)
	case "clear":
		count := list.Clear()
		if err := st.Save(list); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "cleared %d tasks\n", count)
	}
	return nil
}

func parseID(args []string) (int, error) {
	if len(args) < 2 {
		return 0, todo.ErrInvalidID
	}
	id, err := strconv.Atoi(args[1])
	if err != nil || id <= 0 {
		return 0, todo.ErrInvalidID
	}
	return id, nil
}

// Render formats a todo list for humans.
func Render(list todo.List) string {
	if len(list.Tasks) == 0 {
		return "No todos.\n"
	}
	var b strings.Builder
	for _, task := range list.Tasks {
		mark := " "
		if task.Done {
			mark = "x"
		}
		fmt.Fprintf(&b, "[%s] #%d %s\n", mark, task.ID, task.Title)
	}
	return b.String()
}

// Usage returns CLI help text.
func Usage() string {
	return `Usage: todo <command> [args]

Commands:
  add <title>     Add a new todo
  list            List todos
  done <id>       Mark a todo as done
  delete <id>     Delete a todo
  clear           Delete all todos
  help            Show this help

Data file: set JUST_GO_TODO_FILE to choose where JSON is stored.
`
}
