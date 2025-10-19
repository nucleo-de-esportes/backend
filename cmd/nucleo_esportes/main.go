package main

import (
	"flag"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/nucleo-de-esportes/backend/docs"
	"github.com/nucleo-de-esportes/backend/internal/config"
	"github.com/nucleo-de-esportes/backend/internal/handlers"
	"github.com/nucleo-de-esportes/backend/internal/middleware"
	"github.com/nucleo-de-esportes/backend/internal/repository"
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
	// Flag para carregar .env

	vars := flag.String("vars", "file", "Define origem das variáveis de ambiente: 'file' (.env) ou 'exported' (sistema).")
	flag.Parse()

	// Se vars igual a "file", carrega o .env
	if *vars == "file" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Carregar configurações
	cfg := config.LoadConfig()

	repository.Init(cfg.DB)

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
	aulaRoutes := router.Group("/aulas")

	cadRoutes.GET("/mod", handlers.GetAllModalidades)

	cadRoutes.GET("local", handlers.GetAllLocais)

	// // Rotas de turma, sem middleware nem auth
	turmaRoutes.POST("", handlers.CreateTurma)

	turmaRoutes.DELETE("/:id", handlers.DeleteTurma)

	turmaRoutes.GET("/:id", handlers.GetTurmaById)

	turmaRoutes.GET("", handlers.GetAllTurmas)

	turmaRoutes.GET("/nextclass/:id", handlers.GetNextClassById)

	turmaRoutes.PUT("/:id", handlers.UpdateTurma)

	// Rotas de usuário - temporariamente comentadas para testes
	userRoutes := router.Group("/user")
	userRoutes.POST("/register", handlers.RegisterUser)
	userRoutes.GET("", handlers.GetUsers)
	userRoutes.GET("/:id", handlers.GetUserById)
	userRoutes.POST("/login", handlers.LoginUser)
	userRoutes.POST("/inscricao", middleware.AuthUser, handlers.InscreverAluno)
	userRoutes.GET("/turmas", middleware.AuthUser, handlers.GetTurmasByUser)
	userRoutes.DELETE("/:user_id/turma/:turma_id", middleware.AuthUser, handlers.DeleteUserTurma)
	userRoutes.DELETE(("/delete/:id"), middleware.AuthUser, handlers.DeleteUserById)

	aulaRoutes.PUT("/:id/presenca", middleware.AuthUser, handlers.ConfirmarPresenca)

	router.Run(":" + cfg.Server.Port)
}
