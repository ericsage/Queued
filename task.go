package main

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID            int           `json:"id"`
	Progress      float64       `json:"progress"`
	ResourceLocation      string        `json:"resource_location"`
	CreatedOn     time.Time     `json:"created_on"`
	ExpiresOn     time.Time     `json:"expires_on"`
	ProgressCheck ProgressCheck `json:"progress_check"`
}

type ProgressCheck struct {
	endpoint string `json:"endpoint"`
	interval int    `json:"interval"`
}

//Creates a new task and saves it in the queue, which will generate an id for it.
func (q *Queue) NewTask(progress float64, location string, expires time.Time) Task {
	t := Task{}
	t.setTime(expires)
	t.updateLocation(location)
	t.updateProgress(progress)
	q.insertNewTask(&t)
	return t
}

//Update a task that already exists with an id.
func (q *Queue) UpdateTask(id int, progress float64, location string, expires time.Time) Task {
	t := q.readTask(id)
	t.updateExpiry(expires)
	t.updateLocation(location)
	t.updateProgress(progress)
	q.insertTask(&t)
	return t
}

//Set the initial creation time and expiration time. If no expiration time is give, set it to one day.
func (t *Task) setTime(expires time.Time) {
	t.CreatedOn = time.Now()
	if expires.IsZero() {
		expires = time.Now().Add(time.Hour*24)
	}
	t.ExpiresOn = expires
}

func (t *Task) updateExpiry(newExpiry time.Time) {
	if !newExpiry.IsZero() {
		if newExpiry.Before(time.Now()) {
			panic("Cannot set expiry to time in the past")
		}
		t.ExpiresOn = newExpiry
	}
}

/*
Updates the progress of the task. A task cannot be completed until
it's location is set to the created resource. This function should throw
an error if this constraint is not met.
*/
func (t *Task) updateProgress(newProgress float64) {
	if t.ResourceLocation == "" && (t.Progress + newProgress) >= 1 {
		panic("Cannot set progress to 100% unless a location is set for the completed resource")
	}
	t.Progress = t.Progress + newProgress
}

//Updates the location of the tasks created resource.
func (t *Task) updateLocation(newLocation string) {
	if newLocation != "" {
		t.ResourceLocation = newLocation
	}
}

//Converts the task to a json byte string.
func (t *Task) toJSON() []byte {
	buf, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return buf
}
