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

	routes.GET("/:slug/reactions", r.GetPostReactions)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	
	privateRoutes.POST("/", r.CreatePost)
	privateRoutes.PUT("/:slug", r.UpdatePostBySlug)
	privateRoutes.DELETE("/:slug", r.DeletePostBySlug)

	privateRoutes.POST("/:slug/reactions/", r.CreatePostReaction)
	privateRoutes.PUT("/:slug/reactions/:id", r.UpdatePostReaction)
	privateRoutes.DELETE("/:slug/reactions/:id", r.DeletePostReaction)
}
