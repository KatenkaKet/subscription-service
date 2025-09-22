package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

		r.GET("/summary", app.getRecordsByFilter)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(http.StatusFound, "/swagger/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}
