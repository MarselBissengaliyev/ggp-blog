package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/middlewares"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterAuthStoreRoutes(rg *gin.RouterGroup, r *repositories.Repository) {
	routes := rg.Group("/auth")

	routes.POST("/register", r.SignUp)
	routes.POST("/login", r.Login)
	routes.GET("/confirm-email/:email_code", r.ConfirmEmail)

	privateRoutes := routes.Group("/", middlewares.Authenticate(r))
	privateRoutes.GET("/me", r.GetMe)
	privateRoutes.POST("/logout", r.Logout)
}
