package main

import (
  "fmt"
  "os"
  "gopkg.in/yaml.v3"
)

func errNotInMap(task TaskName) error {
  return fmt.Errorf("error: the task %s is not in the list of defined tasks.", task)
}

func errCyclicDep(task TaskName) error {
  return fmt.Errorf("error: the task %s is cyclic.", task)
}

type Action string
type TaskName string

type Task struct {
  Actions *[]Action
  Dependencies *[]TaskName
}

type Digger interface {
  Dig(int) 

  Fill(int, *[][]TaskName)
}

type Tree struct {
  head Digger
  layers [][]TaskName
}

func(t *Tree) Dig() {
  t.head.Dig(0)
}

func(t *Tree) Fill() {
  t.layers = make([][]TaskName, 0)
  t.head.Fill(0, &t.layers)
}

func(t *Tree) PrintLayers() {
  if t.layers == nil {
    fmt.Println("Layers are empty")
    return
  }

  for i, layer := range t.layers {
    fmt.Printf("Layer %d\n", i)
    for _, name:= range layer {
      fmt.Printf("%s ", name)
    }

    fmt.Println()
  }
}

func(t *Tree) Build(name TaskName, tasks map[TaskName]*Task) (err error) {
  t.head, err = buildTree(name, tasks, name, false) 
  return
}

type Leaf struct {
  TaskName
}

type Node struct {
  TaskName
  Children []Digger
}

func (l *Leaf) Dig(depth int) {
  fmt.Printf("depth %d, task %s needs to be run\n", depth, l.TaskName)
}

func (l *Leaf) Fill(depth int, out *[][]TaskName) {
  length := len(*out) 

  fmt.Printf("Length of array is %d, we are at depth %d\n", length - 1,depth)
  for i := length; i < depth + 1; i++ {
    (*out) = append(*out, make([]TaskName, 0))
  }

  (*out)[depth] = append((*out)[depth], l.TaskName)
}

func (n *Node) Dig(depth int) {
  fmt.Printf("depth %d, task %s needs to be run\n", depth, n.TaskName)

  for _, c := range n.Children {
    c.Dig(depth + 1)
  }
}

func (n *Node) Fill(depth int, out *[][]TaskName) {
  length := len(*out)

  fmt.Printf("Length of array is %d, we are at depth %d\n", length - 1,depth)
  for i := length; i < depth + 1; i++ {
    (*out) = append(*out, make([]TaskName, 0))
  }

  (*out)[depth] = append((*out)[depth], n.TaskName) 

  for _, c := range n.Children {
    c.Fill(depth + 1, out)
  }
}

func buildTree(name TaskName, tasks map[TaskName]*Task, initial TaskName, checkCycle bool) (Digger, error) {

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

  children := make([]Digger, len(*task.Dependencies))

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

  var tree Tree
  err = tree.Build(taskName, tasks)

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  tree.Dig()

  tree.Fill()

  tree.PrintLayers()
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}
