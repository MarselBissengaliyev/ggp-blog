package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterUserStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/users")

	routes.GET("/", r.GetUsers)
	routes.GET("/:user_name", r.GetUserByUsername)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	privateRoutes.PUT("/:user_name", r.UpdateUserByUsername)
	privateRoutes.DELETE("/:user_name", r.DeleteUserByUsername)
}
