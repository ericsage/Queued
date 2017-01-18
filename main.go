package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setRoutes(q *Queue, r *gin.Engine) {

	//Return 202 Accepted with a Location header pointing to the task in the queue.
	r.POST("/queue", func(c *gin.Context) {
		defer catchError(c)
		var task Task
		if c.BindJSON(&task) == nil {
			t := q.NewTask(task.Progress, task.ResourceLocation, task.ExpiresOn)
			c.Header("Location", fmt.Sprintf("/queue/%v", t.ID))
			c.IndentedJSON(202, t)
		}
	})

	//Redirect to the resource if the Task is completed, otherwise return Task info.
	r.GET("/queue/:task_id", func(c *gin.Context) {
		defer catchError(c)
		tid, _ := strconv.Atoi(c.Param("task_id"))
		t := q.readTask(tid)
		if t.Progress >= 1.0 {
			c.Redirect(http.StatusMovedPermanently, t.ResourceLocation)
		} else {
			c.IndentedJSON(200, gin.H{
				"percentComplete": (t.Progress * 100),
				"created_on":      t.CreatedOn,
				"expires_on":      t.ExpiresOn,
			})
		}
	})

	//Update the Task with new information
	r.PUT("queue/:task_id", func(c *gin.Context) {
		defer catchError(c)
		tid, _ := strconv.Atoi(c.Param("task_id"))
		var task Task
		if c.BindJSON(&task) == nil {
			t := q.UpdateTask(tid, task.Progress, task.ResourceLocation, task.ExpiresOn)
			c.IndentedJSON(200, t)
		}
	})
}

func catchError(c *gin.Context) {
	if r := recover(); r != nil {
		c.IndentedJSON(400, gin.H{
			"error": r,
		})
	}
}

func main() {
	q := OpenQueue()
	defer q.Close()
	r := gin.Default()
	setRoutes(&q, r)
	r.Run("0.0.0.0:80")
}
