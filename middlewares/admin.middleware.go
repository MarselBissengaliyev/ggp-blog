package middlewares

import (
	"fmt"
	"net/http"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func IsAdmin(r *repositories.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		uid := fmt.Sprint(c.Keys["uid"])

		if uid == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"error":   "uid cannot be emtpy string",
				"message": "error occured while reading uid key",
			})

			c.Abort()

			return
		}

		if err := r.DB.First(&user, fmt.Sprintf("id = %s", uid)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"error":   err.Error(),
				"message": "error occured while finding user by uid",
			})

			c.Abort()

			return
		}

		if user.Role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "failed",
				"error":   "you don't have permission to access this",
				"message": "error occured while checking user role",
			})

			c.Abort()

			return
		}

		c.Next()
	}
}
