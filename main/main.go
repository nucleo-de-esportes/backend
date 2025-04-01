package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/config"
	"github.com/nucleo-de-esportes/backend/controller"
)

func main() {

	supbaseClient := config.InitSupabase()

	router := gin.Default()

	turmaRoutes := router.Group("/turmas")

	turmaRoutes.POST("", func(c *gin.Context) {
		controller.CreateTurma(c, supbaseClient)
	})

	turmaRoutes.DELETE("/:id", func(c *gin.Context) {
		controller.DeleteTurma(c, supbaseClient)

	})

	turmaRoutes.GET("/:id", func(c *gin.Context) {
		controller.ViewTurma(c, supbaseClient)

	})

	turmaRoutes.GET("", func(c *gin.Context) {
		controller.GetAllTurmas(c, supbaseClient)

	})
	router.Run(":8080")
}
