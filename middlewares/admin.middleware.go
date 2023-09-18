package middlewares

import (
	"fmt"
	"net/http"

	"github.com/MarselBissengaliyev/ggp-blog/constants"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/gin-gonic/gin"
)

func IsAdmin(r *repositories.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := fmt.Sprint(c.Keys["role"])

		if role != constants.ADMIN {
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
