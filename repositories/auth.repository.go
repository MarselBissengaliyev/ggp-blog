package repositories

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/constants"
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
		return
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
		return
	}

	r.DB.Model(&models.User{}).Where("user_name = ?", user.UserName).Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "this username already exists",
			"message": "error occured while creating account",
		})
		return
	}

	user.Role = constants.USER
	user.EmailConfirmed = false

	code := utils.Encode(randstr.String(20))
	user.EmailCode = &code
	user.IsBanned = false

	if err := r.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating account",
		})
		return
	}

	emailData := utils.EmailData{
		URL:       r.Config.App_Origin + "/auth/confirm-email/" + *user.EmailCode,
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
		return
	}

	if err := r.DB.Where(&models.User{Email: user.Email}).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "email or password is incorrect",
		})
		return
	}

	passwordIsValid, msg := utils.VerifyPassword(foundUser.Password, user.Password)

	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"error":   "error occured while login",
			"message": msg,
		})
		return
	}

	if !foundUser.EmailConfirmed {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"error":   "email not verificated",
			"message": "Please verify your email",
		})
		c.Abort()
		return
	}

	token, refreshToken, generateTokenErr := utils.GenerateAllTokens(
		foundUser.Email,
		foundUser.LastName,
		foundUser.LastName,
		foundUser.UserName,
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
		return
	}

	newToken, err := utils.CreateTokens(
		token,
		refreshToken,
		c.Request.UserAgent(),
		r.DB,
		foundUser.ID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while create tokens in database",
		})
		return
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
	userId, err := strconv.ParseUint(fmt.Sprint(c.Keys["uid"]), 10, 32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while convert tokenId to uint",
		})
		return
	}

	parsedUserId := uint(userId)

	if err := utils.DeleteTokens(
		r.DB,
		parsedUserId,
		c.Request.UserAgent(),
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "you succefully logout",
	})
}

func (r *Repository) ConfirmEmail(c *gin.Context) {
	var user models.User
	emailCode := c.Param("email_code")

	if emailCode == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"error":   "param email_code cannot be empty value",
			"message": "error occured while checking email value",
		})
		return
	}

	parsedCode := sql.NullString{
		String: emailCode,
		Valid:  true,
	}

	if !parsedCode.Valid {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   "not valid email code",
			"message": "error while validating email_code",
		})
		return
	}

	if err := r.DB.Where("email_code = ?", parsedCode.String).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while trying to found a user with this email_code",
		})

		return
	}

	user.EmailCode = nil

	user.EmailConfirmed = true

	if err := r.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while trying to save a user with this email_code",
		})

		return
	}

	token, refreshToken, err := utils.GenerateAllTokens(
		user.Email,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.ID,
		user.Role,
		r.Config.Token_Secret,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while generate tokens",
		})
		return
	}

	createdToken, err := utils.CreateTokens(token, refreshToken, c.Request.UserAgent(), r.DB, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"token":   createdToken.AccessToken,
		"message": "you succefully confirmed email",
	})
}

func (r *Repository) GetMe(c *gin.Context) {
	var user models.User

	uid := fmt.Sprint(c.Keys["uid"])

	if err := r.DB.First(&user, fmt.Sprintf("id = %s", uid)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting user information",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    user,
		"message": "you successfully got your information",
	})
}
