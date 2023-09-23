package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterTagStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/posts/:post_id/tags")

	routes.GET("/", r.GetTags)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	
	privateRoutes.POST("/", r.CreateTag)
	privateRoutes.PUT("/:tag_id", r.UpdateTag)
	privateRoutes.DELETE("/:tag_id", r.DeleteTag)
}
