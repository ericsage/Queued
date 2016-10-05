package main

import (
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
)

type Queue struct {
	DB *bolt.DB
}

var QUEUE = []byte("QUEUE")

//Open the Queue DB and make sure a bucket exists for Tasks.
func OpenQueue() Queue {
	db, err := bolt.Open("queued.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(QUEUE)
		return err
	})
	return Queue{DB: db}
}

//Generate a unique ID for the Task before inserting it into the Queue.
func (q *Queue) insertNewTask(t *Task) *Task {
	q.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(QUEUE)
		id, _ := b.NextSequence()
		t.ID = int(id)
		return b.Put(itob(t.ID), t.toJSON())
	})
	return t
}

//Insert an updated Task into the Queue.
func (q *Queue) insertTask(t *Task) *Task {
	q.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(QUEUE)
		return b.Put(itob(t.ID), t.toJSON())
	})
	return t
}

//Read the Task from the Queue.
func (q *Queue) readTask(id int) Task {
	t := Task{}
	q.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(QUEUE)
		json.Unmarshal(b.Get(itob(id)), &t)
		return nil
	})
	return t
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//Close the Queue's DB. Should be defered by main.
func (q *Queue) Close() {
	q.DB.Close()
}
