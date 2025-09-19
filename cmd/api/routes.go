package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	r := g.Group("/api/subscriptions")
	{
		r.GET("/all", app.getAllRecords)
		r.GET("/:id", app.getRecordByID)
		r.GET("/user/:id", app.getRecordsByUserID)
		r.GET("/service/:name", app.getRecordsByServiceName)
		r.PUT("/:id", app.updateRecordByID)
		r.DELETE("/:id", app.deleteRecordByID)
		r.DELETE("/user/:id", app.deleteRecordByUserID)
		r.DELETE("/service/:name", app.deleteRecordsByServiceName)
		r.POST("/newrecord", app.createRecord)
	}

	return g
}
