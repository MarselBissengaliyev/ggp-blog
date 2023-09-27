package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterTagStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/posts/:slug/tags")

	routes.GET("/", r.GetTags)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	
	privateRoutes.PUT("/", r.UpdateTags)
	privateRoutes.PUT("/", r.DeleteTag)
}
