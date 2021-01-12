package main

import (
	"fmt"
)

// Filler builds a tree and a stack of depedencies
type Filler interface {
	Fill(int, *[][]TaskName)
}

// Tree contains the tree and the layers
type Tree struct {
	Tasks  map[TaskName]*Task
	head   Filler
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
	t.head, err = t.build(name, name, false)
	if err != nil {
		return
	}

	fmt.Println("Tree has been built, cutting in layers...")

	t.fill()

	fmt.Println("Done!")

	return
}

func (t Tree) build(name TaskName, initial TaskName, checkCycle bool) (Filler, error) {
	if checkCycle && name == initial {
		return nil, errCyclicDep(name)
	}

	task := t.Tasks[name]

	if task == nil {
		return nil, errNotInMap(name)
	}

	if task.Dependencies == nil {
		return &Leaf{TaskName: name}, nil
	}

	children := make([]Filler, len(*task.Dependencies))

	for i, child := range *task.Dependencies {
		digger, err := t.build(child, initial, true)

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

func (t Tree) Run() (err error) {
	length := len(t.layers)

	kill, done := make(chan error), make(chan error)
	defer close(kill)
	defer close(done)

	for length > 0 {
		fmt.Printf("running layer %d\n", length)
		layer := t.layers[length - 1]
		layerLength := len(layer)

		for _, name := range layer {
			fmt.Printf("\tstarting task %s\n", name)
			go t.Tasks[name].Run(done, kill)
		}

		// while some tasks are not done
		for layerLength > 0 {
			if err := <- done; err != nil {
				kill <- err
				return err
			}
			layerLength--
		}

		length--
		fmt.Printf("layer %d done\n", length)
	}

	return nil
}
