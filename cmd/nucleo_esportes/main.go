package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/internal/handlers"
	"github.com/nucleo-de-esportes/backend/internal/repository"

	_ "github.com/nucleo-de-esportes/backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Nucleo de Esportes API
// @version 1.0
// @description API do sistema de gerenciamento de turmas do núcleo de esportes da faculdade.
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func main() {

	repository.Init()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	turmaRoutes := router.Group("/turmas")

	cadRoutes := router.Group("/cad")

	cadRoutes.GET("/mod", handlers.GetAllModalidades)

	cadRoutes.GET("local", handlers.GetAllLocais)

	// Rotas de turma, sem middleware nem auth
	turmaRoutes.POST("", handlers.CreateTurma)

	turmaRoutes.DELETE("/:id", handlers.DeleteTurma)

	turmaRoutes.GET("/:id", handlers.GetTurmaById)

	turmaRoutes.GET("", handlers.GetAllTurmas)

	turmaRoutes.PUT("/:id", handlers.UpdateTurma)

	// Rotas de usuário - temporariamente comentadas para testes
	// userRoutes := router.Group("/user")
	// userRoutes.POST("/register", handlers.RegsiterUser)
	// userRoutes.POST("/login", handlers.LoginUser)
	// userRoutes.POST("/inscricao", handlers.InscreverAluno)
	// userRoutes.GET("/turmas", handlers.GetTurmasByUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	router.Run(port)
}
