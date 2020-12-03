package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

func errNotInMap(task TaskName) error {
	return fmt.Errorf("error: the task %s is not in the list of defined tasks.", task)
}

func errCyclicDep(task TaskName) error {
	return fmt.Errorf("error: the task %s is cyclic.", task)
}

// Action is a command to run
type Action string

// TaskName is the name of a task
type TaskName string

// Task is a list of Actions and an optional list of TaskNames it depends on
type Task struct {
	Actions      *[]Action
	Dependencies *[]TaskName
}

// TreeBuilder builds a tree and a stack of depedencies
type TreeBuilder interface {
	Fill(int, *[][]TaskName)
}

// Tree contains the tree and the layers
type Tree struct {
  Tasks  map[TaskName]*Task
	head   TreeBuilder
	layers [][]TaskName
}

func (t *Tree) removeDuplicates() {
	newLayers := make([][]TaskName, len(t.layers))

	taskMap := map[TaskName]int{}

	for depth, layer := range t.layers {
		for _, name := range layer {
			taskMap[name] = depth
		}
	}

	for task, index := range taskMap {
		newLayers[index] = append(newLayers[index], task)
	}

	t.layers = newLayers
}

// fill fills the layers from the tree representation
func (t *Tree) fill() {
	t.layers = make([][]TaskName, 0)
	t.head.Fill(0, &t.layers)
  fmt.Println("layers have been built, optimising...")
	t.removeDuplicates()
}

func (t *Tree) PrintLayers() {
	if t.layers == nil {
		fmt.Println("Layers are empty")
		return
	}

	for i, layer := range t.layers {
		printer := "Layer %d: [ "

		for _, name := range layer {
			printer += fmt.Sprintf("%s ", name)
		}

		printer += "]\n"

		fmt.Printf(printer, i)
	}
}

func (t *Tree) Build(name TaskName) (err error) {
	t.head, err = buildTree(name, t.Tasks, name, false)
	if err != nil {
		return
	}

  fmt.Println("Tree has been built, cutting in layers...")

	t.fill()

  fmt.Println("Done!")

  return
}

func (t *Tree) Run() (err error) {
  length := len(t.layers)

  for length > 0 {
    for _, name := range t.layers[length-1] {
      task := t.Tasks[name]

      for _, action := range *task.Actions {
        fmt.Printf("%s: %s\n", name, action)
        strCmd := strings.Split(string(action), " ")

        command := strCmd[0]
        args := strCmd[1:]

        cmd := exec.Command(command, args...)
        cmd.Stdout = os.Stdout

        err := cmd.Run()

        if err != nil {
          return err
        }
      }
    }

    length--
  }

  return nil
}

// Leaf represents a task without dependencies
type Leaf struct {
	TaskName
}

// Node represents a task and it's dependencies
type Node struct {
	TaskName
	Children []TreeBuilder
}

func (l *Leaf) Fill(depth int, out *[][]TaskName) {
	length := len(*out)

	for i := length; i < depth+1; i++ {
		(*out) = append(*out, make([]TaskName, 0))
	}

	(*out)[depth] = append((*out)[depth], l.TaskName)
}

func (n *Node) Fill(depth int, out *[][]TaskName) {
	length := len(*out)

	for i := length; i < depth+1; i++ {
		(*out) = append(*out, make([]TaskName, 0))
	}

	(*out)[depth] = append((*out)[depth], n.TaskName)

	for _, c := range n.Children {
		c.Fill(depth+1, out)
	}
}

func buildTree(name TaskName, tasks map[TaskName]*Task, initial TaskName, checkCycle bool) (TreeBuilder, error) {

	if checkCycle && name == initial {
		return nil, errCyclicDep(name)
	}

	task := tasks[name]

	if task == nil {
		return nil, errNotInMap(name)
	}

	if task.Dependencies == nil {
		return &Leaf{TaskName: name}, nil
	}

	children := make([]TreeBuilder, len(*task.Dependencies))

	for i, child := range *task.Dependencies {
		digger, err := buildTree(child, tasks, initial, true)

		if err != nil {
			return nil, err
		}

		children[i] = digger
	}

	return &Node{
		TaskName: name,
		Children: children,
	}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("runnerup requires exactly 1 argument.")
		fmt.Println("usage :")
		fmt.Println("\trunnerup <task>")
		return
	}

	taskName := TaskName(os.Args[1])

	file, err := os.Open("tasks.yaml")
	check(err)

	var tasks map[TaskName]*Task
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&tasks)
	check(err)

  tree := Tree{Tasks: tasks}
	err = tree.Build(taskName)
	check(err)

  fmt.Println()

  err = tree.Run()
  check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
