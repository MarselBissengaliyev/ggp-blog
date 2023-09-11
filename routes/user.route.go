package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserStoreRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	routes := rg.Group("/users")

	routes.GET("/")
	routes.GET("/:username")
}
