package main

import (
	"encoding/json"
	"log"
)

type Task struct {
	ID       int `json:"id"`
	Progress float64 `json:"progress"`
	Location string
}

//Creates a new task and saves it in the queue, which will generate an id for it.
func (q *Queue) NewTask(progress float64, location string) Task {
	t := Task{}
	t.updateLocation(location)
	t.updateProgress(progress)
	q.insertNewTask(&t)
	return t
}

//Update a task that already exists with an id.
func (t * Task) Update(progress float64, location string) {
	t.updateLocation(location)
	t.updateProgress(progress)
	q.insertTask(&t)
}

/*
Updates the progress of the task. A task cannot be completed until
it's location is set to the created resource. This function should throw
an error if this constraint is not met.
*/
func (t *Task) updateProgress(newProgress float64) {
	if t.Location == "" && newProgress == 1 {
		panic("Cannot set progress to 100% unless a location is set")
	}
	t.Progress = newProgress
}

//Updates the location of the tasks created resource.
func (t *Task) updateLocation(newLocation string) {
	t.Location = newLocation
}

//Converts the task to a json byte string.
func (t *Task) toJSON() []byte {
	buf, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}
