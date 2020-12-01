package main

import (
  "fmt"
  "os"
  "gopkg.in/yaml.v3"
)

type Action string
type TaskName string

type Task struct {
  Actions *[]Action
  Dependencies *[]TaskName
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

  task := tasks[taskName]

  if task == nil {
    fmt.Printf("The task %s was not found in the file tasks.yaml\n", taskName)
    return
  }

  fmt.Printf("The task %s was found in tasks.yaml!!!\n details: \nActions: %+v\nDependencies: %+v\n", taskName, *task.Actions, *task.Dependencies)
  return
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}
