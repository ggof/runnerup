package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	length := len(os.Args)
	if length != 2 {
		fmt.Println(errOneArgument(length - 1))
		fmt.Println(getHelp())
		return
	}

	taskName := TaskName(os.Args[1])

	file, err := os.Open("tasks.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var tasks map[TaskName]*Task
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

  tree := Tree{Tasks: tasks}
	err = tree.Build(taskName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

  fmt.Println()

  err = tree.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func getHelp() string {
	usage := "usage: runnerup <task>\n"
	usage += "\ttask :\ta task that can be found in the file \"tasks.yaml\"\n"
	usage += "\t\ttype: string"

	return usage
}
