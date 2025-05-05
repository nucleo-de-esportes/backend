package main

import (
	"flag"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/config"
	"github.com/nucleo-de-esportes/backend/controller"
)

func main() {
	env := flag.String("vars", "file", "Defines from where to load env vars: file or exported")

	flag.Parse()

	supbaseClient := config.InitSupabase(*env)

	router := gin.Default()

	router.Use(cors.Default())

	turmaRoutes := router.Group("/turmas")

	userRoutes := router.Group("/user")

	turmaRoutes.POST("", func(c *gin.Context) {
		controller.CreateTurma(c, supbaseClient)
	})

	turmaRoutes.DELETE("/:id", func(c *gin.Context) {
		controller.DeleteTurma(c, supbaseClient)

	})

	turmaRoutes.GET("/:id", func(c *gin.Context) {
		controller.GetTurmaById(c, supbaseClient)

	})

	turmaRoutes.GET("", func(c *gin.Context) {
		controller.GetAllTurmas(c, supbaseClient)

	})

	turmaRoutes.PUT("/:id", func(c *gin.Context) {
		controller.UpdateTurma(c, supbaseClient)
	})

	userRoutes.POST("/register", func(c *gin.Context) {
		controller.RegsiterUser(c, supbaseClient)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	router.Run(port)
}
