package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterPostReactionStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/posts/:slug/reactions")

	routes.GET("/", r.GetReactions)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))

	privateRoutes.POST("/", r.CreateReaction)
	privateRoutes.PUT("/:reaction_id", r.UpdateReaction)
	privateRoutes.DELETE("/:reaction_id", r.DeleteReaction)
}
