package routes

import (
	"github.com/DevNCH/NAS/internal/database"
	"github.com/DevNCH/NAS/internal/repository"
	"github.com/DevNCH/NAS/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	db := database.GetDB()

	userRepo := repository.NewUserRepository(db)

	authService := services.NewAuthService(userRepo)

	_ = authService // temporário para não dar erro de variável não utilizada

	router.LoadHTMLGlob("web/templates/*")

	router.Static("/static", "./web/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Servidor NAS",
		})
	})

	return router
}
