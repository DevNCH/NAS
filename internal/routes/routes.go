package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

    router := gin.Default()

	router.GET("/", func(c *gin.Context) {
    	c.String(200, "Servidor NAS funcionando!")
	})

    return router
}