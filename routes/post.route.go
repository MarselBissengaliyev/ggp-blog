package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterPostStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/posts")

	routes.GET("/", r.GetPosts)
	routes.GET("/:slug", r.GetPostBySlug)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	
	privateRoutes.POST("/", r.CreatePost)
	privateRoutes.PUT("/:slug", r.UpdatePostBySlug)
	privateRoutes.DELETE("/:slug", r.DeletePostBySlug)
}
