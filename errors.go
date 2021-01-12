package main

import "fmt"

func errOneArgument(length int) error {
	return fmt.Errorf("error: one argument expected, but given %d", length)
}

func errNotInMap(task TaskName) error {
	return fmt.Errorf("error: the task %s is not in the list of defined tasks", task)
}

func errCyclicDep(task TaskName) error {
	return fmt.Errorf("error: the task %s is cyclic", task)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
