package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupUsersRoutes(router *gin.Engine) {
	router.GET("/users", getUsers)
}

func getUsers(c *gin.Context) {
	// Lógica para obtener datos de la colección 'points'
}
