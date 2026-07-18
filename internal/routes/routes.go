package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

    router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")

	router.Static("/static", "./web/static")

	router.GET("/", func(c *gin.Context) {
    	c.HTML(200, "index.html", gin.H{
			"title": "Servidor NAS",
		})
	})

    return router
}