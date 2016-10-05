package main

import "github.com/gin-gonic/gin"

func setRoutes(q Queue, r int) {

  //Return 202 Accepted with a Location header pointing to the task in the queue.
  r.POST("/queue", func(c *gin.Context) {
    t := q.NewTask(0.1, "http://example.com/myresource/43424")
    c.Header("Location", t.ID)
    c.JSON(202, gin.H{
      "task_id": t.ID
    })
  })

  //Redirect to the resource if the Task is completed, otherwise return Task info.
  r.GET("/queue/:task_id", func(c *gin.Context) {
    tid := c.Param("task_id")
    t := q.readTask(tid)
    if (t.Progress == 1.0) {
      c.Redirect(http.StatusMovedPermanently, t.Location)
    } else {
      c.JSON(200, gin.H{
        "progress": t.Progress
      })
    }
  })

  //Update the Task with new information
  r.PUT("queue/:task_id", func(c *gin.Context) {
    tid := c.Param("task_id")
    t := q.readTask(tid)
  })
}

func main() {
	q := OpenQueue()
  defer q.Close()
  r := gin.Default()
  setRoutes(q, r)
	r.Run()
	fmt.Println(t1.Location)
}
