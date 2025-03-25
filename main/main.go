package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/config"
	"github.com/nucleo-de-esportes/backend/controller"
)

func main() {

	supbaseClient := config.InitSupabase()

	request := gin.Default()

	request.POST("/turma", func(c *gin.Context) {
		controller.CreateTurma(c, supbaseClient)
	})

	request.Run(":8080")
}
