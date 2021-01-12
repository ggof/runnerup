package main

import (
	"os"
	"os/exec"
	"strings"
)

// Action is a command to run
type Action string

// TaskName is the name of a task
type TaskName string

// Task is a list of Actions and an optional list of TaskNames it depends on
type Task struct {
	Actions      *[]Action
	Dependencies *[]TaskName
}

// Leaf represents a task without dependencies
type Leaf struct {
	TaskName
}

// Node represents a task and it's dependencies
type Node struct {
	TaskName
	Children []Filler
}

func (l *Leaf) Fill(depth int, out *[][]TaskName) {
	length := len(*out)

	for i := length; i < depth+1; i++ {
		*out = append(*out, make([]TaskName, 0))
	}

	(*out)[depth] = append((*out)[depth], l.TaskName)
}

func (n *Node) Fill(depth int, out *[][]TaskName) {
	length := len(*out)

	for i := length; i < depth+1; i++ {
		*out = append(*out, make([]TaskName, 0))
	}

	(*out)[depth] = append((*out)[depth], n.TaskName)

	for _, c := range n.Children {
		c.Fill(depth+1, out)
	}
}

func (t Task) Run(sendError chan<- error, recvError <-chan error) {
	done, kill := make(chan error), make(chan struct{})
	defer close(done)
	defer close(kill)

	for _, action := range *t.Actions {

		go action.Run(done, kill)

		select {
		case err := <-done:
			if err != nil {
				sendError <- err
				return
			}
		case <-recvError:
			kill <- struct{}{}
			return
		}
	}

	// once all of this is done, send nil up (this task is done)
	sendError <- nil
}

func (a *Action) Run(done chan error, kill chan struct{}) {
	strCmd := strings.Split(string(*a), " ")

	command := strCmd[0]
	args := strCmd[1:]

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout

	select {
	case <-kill:
		_ = cmd.Process.Kill()
	case done <- cmd.Run():
	}
}

