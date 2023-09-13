package middlewares

import (
	"net/http"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/MarselBissengaliyev/ggp-blog/repositories"
	"github.com/MarselBissengaliyev/ggp-blog/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(r *repositories.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"error":   "No authorization header provided",
				"message": "error occured while getting token from header",
			})

			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(clientToken, r.Config)

		if err != "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "failed",
				"error":   err,
				"message": "error occured while validate token",
			})
			c.Abort()
			return
		}

		if err := r.DB.First(&user, claims.Uid).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "failed",
				"error":   err.Error(),
				"message": "the user belonging to this token no longer exists",
			})
			c.Abort()
			return
		}

		if err := r.DB.First(&models.Token{
			UserId:    user.ID,
			UserAgent: c.Request.UserAgent(),
		}).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "failed",
				"error":   err.Error(),
				"message": "the token belonging to this user no longer exists",
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("role", claims.Role)
		c.Next()
	}
}
