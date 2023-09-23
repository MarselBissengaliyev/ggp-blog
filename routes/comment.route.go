package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterCommentStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/posts/:post_id/comments")

	routes.GET("/", r.GetComments)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))

	privateRoutes.POST("/", r.CreateComment)
	privateRoutes.PUT("/:id", r.UpdateComment)
	privateRoutes.DELETE("/:id", r.DeleteComment)
}
