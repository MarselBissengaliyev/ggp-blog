package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterPostStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/posts")

	routes.GET("/", r.GetPosts)
	routes.GET("/:post_id", r.GetPostById)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	
	privateRoutes.POST("/", r.CreatePost)
	privateRoutes.PUT("/:post_id", r.UpdatePostById)
	privateRoutes.DELETE("/:post_id", r.DeletePostById)
}
