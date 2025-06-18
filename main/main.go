package main

import (
	"flag"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/config"
	"github.com/nucleo-de-esportes/backend/controller"
	"github.com/nucleo-de-esportes/backend/main/middleware"

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
func main() {
	env := flag.String("vars", "file", "Defines from where to load env vars: file or exported")

	flag.Parse()

	supbaseClient := config.InitSupabase(*env)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	turmaRoutes := router.Group("/turmas")

	userRoutes := router.Group("/user")

	cadRoutes := router.Group("/cad")

	cadRoutes.GET("/mod", func(c *gin.Context) {
		controller.GetAllModalidades(c, supbaseClient)
	})

	cadRoutes.GET("local", func(c *gin.Context) {
		controller.GetAllLocais(c, supbaseClient)
	})

	turmaRoutes.POST("", middleware.AuthUser(supbaseClient), func(c *gin.Context) {
		controller.CreateTurma(c, supbaseClient)
	})

	turmaRoutes.DELETE("/:id", middleware.AuthUser(supbaseClient), func(c *gin.Context) {
		controller.DeleteTurma(c, supbaseClient)

	})

	turmaRoutes.GET("/:id", middleware.AuthUser(supbaseClient), func(c *gin.Context) {
		controller.GetTurmaById(c, supbaseClient)

	})

	turmaRoutes.GET("", func(c *gin.Context) {
		controller.GetAllTurmas(c, supbaseClient)

	})

	turmaRoutes.PUT("/:id", middleware.AuthUser(supbaseClient), func(c *gin.Context) {
		controller.UpdateTurma(c, supbaseClient)
	})

	userRoutes.POST("/register", func(c *gin.Context) {
		controller.RegsiterUser(c, supbaseClient)
	})

	userRoutes.POST("/login", func(c *gin.Context) {
		controller.LoginUser(c, supbaseClient)
	})

	userRoutes.POST("/inscricao", middleware.AuthUser(supbaseClient), func(c *gin.Context) {
		controller.InscreverAluno(c, supbaseClient)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	router.Run(port)
}
