package routes

import (
	"time"

	"github.com/DevNCH/NAS/internal/auth"
	"github.com/DevNCH/NAS/internal/database"
	"github.com/DevNCH/NAS/internal/handlers"
	"github.com/DevNCH/NAS/internal/middleware"
	"github.com/DevNCH/NAS/internal/repository"
	"github.com/DevNCH/NAS/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	db := database.GetDB()

	userRepo := repository.NewUserRepository(db)

	authService := services.NewAuthService(userRepo)

	sessions := auth.NewSessionManager(24 * time.Hour)

	authHandler := handlers.NewAuthHandler(authService, sessions)

	fileRepo := repository.NewFileRepository(db)

	fileService := services.NewFileService(fileRepo)

	fileHandler := handlers.NewFileHandler(fileService)

	router.LoadHTMLGlob("web/templates/*")

	router.Static("/static", "./web/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Servidor NAS",
		})
	})

	router.GET("/files", fileHandler.ListFiles)

	// Rotas de autenticação (públicas)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	// Rotas protegidas: exigem sessão válida
	protected := router.Group("/")
	protected.Use(middleware.RequireAuth(sessions))
	{
		protected.GET("/me", authHandler.Me)
	}

	// Exemplo de rota restrita a administradores
	admin := router.Group("/admin")
	admin.Use(middleware.RequireAuth(sessions), middleware.RequireRole("admin"))
	{
		admin.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong, admin"})
		})
	}

	return router
}
