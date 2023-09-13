package repositories

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/MarselBissengaliyev/ggp-blog/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

func (r *Repository) SignUp(c *gin.Context) {
	var user models.User
	var count int64

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding json",
		})
	}

	password := utils.HashPassword(user.Password)
	user.Password = password

	r.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "this email already exists",
			"message": "error occured while creating account",
		})
	}

	user.Role = "USER"
	user.EmailConfirmed = false

	code := randstr.String(20)
	user.EmailCode = sql.NullString{String: utils.Encode(code), Valid: true}
	user.IsBanned = false

	if err := r.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating account",
		})
	}

	emailData := utils.EmailData{
		URL:       r.Config.Client_Origin + "/auth/confirm-email/" + user.EmailCode.String,
		FirstName: user.FirstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&user, &emailData, r.Config)

	message := "We sent an email with a verification code to " + user.Email

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
	})
}

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

	if err := r.DB.Where(&models.User{Email: user.Email}).First(&foundUser).Error; err != nil {
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
		r.Config.Token_Secret,
	)

	if generateTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   generateTokenErr.Error(),
			"message": "error occured while generate tokens",
		})
	}

	newToken, err := utils.CreateTokens(
		token,
		refreshToken,
		c.Request.UserAgent(),
		foundUser.ID,
		r.DB,
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

func (r *Repository) Logout(c *gin.Context) {
	tokenId, err := strconv.ParseUint(fmt.Sprint(c.Keys["token_id"]), 10, 32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while convert tokenId to uint",
		})
	}

	parsedTokenId := uint(tokenId)

	if err := utils.DeleteTokens(r.DB, parsedTokenId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete tokens",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "you succefully logout",
	})
}
