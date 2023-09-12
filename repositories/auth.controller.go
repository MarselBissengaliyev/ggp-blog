package repositories

import (
	"net/http"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/MarselBissengaliyev/ggp-blog/utils"
	"github.com/gin-gonic/gin"
)

func (r *Repository) Login(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding json",
		})
	}

	if err := r.db.Where(&models.User{Email: user.Email}).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "email or password is incorrect",
		})
	}

	passwordIsValid, msg := utils.VerifyPassword(user.Password, foundUser.Password)

	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"error":   "error occured while login",
			"message": msg,
		})
	}

	token, refreshToken, generateTokenErr := utils.GenerateAllTokens(
		foundUser.Email,
		foundUser.LastName,
		foundUser.LastName,
		foundUser.ID,
		foundUser.Role,
	)

	if generateTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   generateTokenErr.Error(),
			"message": "error occured while generate tokens",
		})
	}

	err, newToken := utils.CreateTokens(
		token,
		refreshToken,
		c.Request.UserAgent(),
		foundUser.ID,
		r.db,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while create tokens in database",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": newToken.AccessToken,
		},
		"message": "you succefully login",
	})
}
