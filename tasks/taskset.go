package tasks

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type task struct {
	name   string
	label  string
	outDir string
	cmd    *exec.Cmd
	err    error
	outF   string
	errF   string
}

func (t *task) Run() {
	name := fmt.Sprintf("%s_%s", t.name, t.label)
	if t.name == "" {
		name = t.label
	}
	t.outF = fmt.Sprintf("%s/%s.out", t.outDir, name)
	t.errF = fmt.Sprintf("%s/%s.err", t.outDir, name)

	of, err := os.Create(t.outF)
	if err != nil {
		t.err = err
		return
	}
	oe, err := os.Create(t.errF)
	if err != nil {
		t.err = err
		return
	}
	t.cmd.Stdout = of
	t.cmd.Stderr = oe

	t.err = t.cmd.Run()
}

func (t *task) PrintStatus() {
	if t.err == nil {
		fmt.Println(t.message("completed: succeeded"))
	} else {
		fmt.Println(t.message("completed: failed"))
		fmt.Println(t.message("===stdout begin==="))
		dumpFile(t.outF)
		fmt.Println(t.message("===stdout end==="))

		// Do we need to dump this on stderr?
		fmt.Println(t.message("===stderr begin==="))
		dumpFile(t.errF)
		fmt.Println(t.message("===stderr end==="))
	}
}

func (t *task) message(m string) string {
	name := t.name
	if name == "" {
		name = "."
	}
	return fmt.Sprintf("[%s]: package %s %s", t.label, name, m)
}

func dumpFile(f string) {
	in, err := os.Open(f)
	if err != nil {
		return
	}
	defer in.Close()
	if _, err = io.Copy(os.Stdout, in); err != nil {
		return
	}
}

// TaskSet executes a set of commands.
type TaskSet struct {
	Name     string
	outDir   string
	tasks    []*task
	waiter   *sync.WaitGroup
	notifier *sync.WaitGroup
	done     bool
}

// New creates a new task set with the specified name.
func New(name string, outDir string, notifier *sync.WaitGroup) *TaskSet {
	return &TaskSet{
		Name:     name,
		outDir:   outDir,
		tasks:    make([]*task, 0, 16),
		waiter:   &sync.WaitGroup{},
		notifier: notifier,
		done:     false,
	}
}

// Add adds a command to this task set for execution.
func (ts *TaskSet) Add(name string, cmd *exec.Cmd) {
	t := &task{name: name, label: ts.Name, outDir: ts.outDir, cmd: cmd}
	ts.tasks = append(ts.tasks, t)
	ts.waiter.Add(1)
}

// Run executes all the commands in parallel. When all commands are done
// it notifies the registered waitgroup. Run returns error is any of the
// command failed else nil. The caller should check individual command status
// for errors if Run returns a failure.
func (ts *TaskSet) Run() error {
	if ts.done {
		// Tasks already executed
		return errors.New("taskset already executed")
	}
	defer func() {
		ts.done = true
		if ts.notifier != nil {
			ts.notifier.Done()
		}
	}()

	// Execute all the tasks and wait
	for _, cmd := range ts.tasks {
		go func(t *task, w *sync.WaitGroup) {
			t.Run()
			w.Done()
		}(cmd, ts.waiter)
	}

	// Wait
	ts.waiter.Wait()

	// Check for error
	for _, t := range ts.tasks {
		if !t.cmd.ProcessState.Success() {
			return errors.New("taskset failed")
		}
	}
	return nil
}

// Done returns true if the task set completed else false.
func (ts *TaskSet) Done() bool {
	return ts.done
}

// PrintStatus prints the status of all the tasks to stdout.
func (ts *TaskSet) PrintStatus() {
	for _, t := range ts.tasks {
		t.PrintStatus()
	}
}
