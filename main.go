package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

const path = "taskList.json"

var taskList = make(map[int]Task)

type Status int

const (
	Planned = iota
	InProgress
	Finished
)

func (s Status) String() string {
	switch s {
	case Planned:
		return "Planned"
	case InProgress:
		return "In Progress"
	case Finished:
		return "Finished"
	default:
		return "Unknown"
	}
}

type Task struct {
	Title       string
	Description string
	Status      Status
}

func main() {
	// os.Args[0] is always the path of exec
	if len(os.Args) < 2 {
		fmt.Println("commands: \n- add {title} {description} {status} \n- list\n- edit {id} {key} {value}")
		return
	}

	load(path)

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		if len(args) != 2 {
			fmt.Println("Wrong number of args (!=2)")
			return
		}
		add(args)
	case "list", "ls":
		list()
		return // no need to save
	case "edit":
		if len(args) != 3 {
			fmt.Println("Wrong number of args (!=3)")
		}
		edit(args)
	case "del", "delete", "rem", "remove":
		if len(args) != 1 {
			fmt.Println("Wrong number of args (!=1)")
		}
		del(args)
	default:
		fmt.Println("Unknown command")
		return
	}
	save()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func save() (err error) {
	encoded, err := json.Marshal(taskList)
	check(err)

	err = os.WriteFile(path, encoded, 0644)
	check(err)
	return
}

func load(path string) (err error) {
	// open file with read/write, if doesn't exist create
	data, err := os.ReadFile(path)
	if !os.IsNotExist(err) {
		check(err)
	}

	// if the file has content try parsing it
	if len(data) > 0 {
		err = json.Unmarshal(data, &taskList)
		check(err)
	}

	return
}

func add(args []string) {
	t := Task{args[0], args[1], Planned}
	taskList[len(taskList)+1] = t
}

func list() {
	for i, t := range taskList {
		fmt.Printf("ID: %d\n- Title: %s \n- Description: %s \n- Status: %s \n", i, t.Title, t.Description, t.Status.String())
	}
}

func edit(args []string) {
	// 0: id, 1: key, 2: value
	id, err := strconv.Atoi(args[0])
	check(err)

	t, ok := taskList[id]
	if !ok {
		fmt.Println("Coulnd't find task")
	}

	switch args[1] {
	case "title", "Title":
		t.Title = args[2]
	case "desc", "Description":
		t.Description = args[2]
	case "Status", "status", "stat":
		switch args[2] {
		case "done", "finished", "fin":
			t.Status = Finished
		case "plan", "todo":
			t.Status = Planned
		case "doing", "wip", "progres":
			t.Status = InProgress
		default:
			fmt.Println("Unknown status type. available: todo, wip, done")
		}
	default:
		fmt.Println("Couldn't find key. available: title, desc, status")
	}
	taskList[id] = t
}

func del(args []string) {
	id, err := strconv.Atoi(args[0])
	check(err)

	_, ok := taskList[id]
	if !ok {
		fmt.Println("Coulnd't find task")
	}

	delete(taskList, id)
}
