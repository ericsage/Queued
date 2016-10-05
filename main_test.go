
package main

import (
  "testing"
)

func TestOpenAndCloseQueue(t *testing.T) {
  q := OpenQueue()
  q.Close()
}

func TestCreateTask(t *testing.T) {
  q := OpenQueue()
  q.NewTask(0, "")
  q.Close()
}

func ReadTask(t *testing.T) {
  q := OpenQueue()
  task := q.NewTask(0, "")
  q.readTask(task.ID)
  q.Close()
}
