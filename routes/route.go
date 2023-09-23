package routes

import (
	"github.com/MarselBissengaliyev/ggp-blog/config"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, config *config.Config) {
	v1 := r.Group("/api/v1")

	repository := &repositories.Repository{
		DB:     db,
		Config: config,
	}

	RegisterUserStoreRoutes(v1, repository)
	RegisterAuthStoreRoutes(v1, repository)
	RegisterPostStoreRoutes(v1, repository)
	RegisterPostReactionStoreRoutes(v1, repository)
	RegisterCommentStoreRoutes(v1, repository)
	RegisterTagStoreRoutes(v1, repository)
}
