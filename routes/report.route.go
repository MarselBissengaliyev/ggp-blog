package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterReportStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/reports")

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	privateRoutes.POST("/", r.CreateReport)

	adminRoutes := privateRoutes.Group("/", middlewares.IsAdmin(r))

	adminRoutes.GET("/", r.GetReports)
	adminRoutes.GET("/:id", r.GetReportById)
}
