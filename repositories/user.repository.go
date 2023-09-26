package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/MarselBissengaliyev/ggp-blog/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

func (r *Repository) GetUsers(c *gin.Context) {
	var users []models.User

	limit := 5

	page, err := strconv.Atoi(c.Query("page"))
	offset := limit * page

	if err != nil || page == 1 {
		page = 1
		offset = 0
	}

	orderBy := c.Query("order_by")

	if orderBy == "" {
		orderBy = "email"
	}

	orderType := c.Query("order_type")

	if orderType == "" {
		orderType = "desc"
	}

	if err := r.DB.Limit(limit).Offset(offset).Order(fmt.Sprintf(
		"%s %s",
		orderBy,
		orderType,
	)).Find(&users).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting users",
		})

		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "success",
		"data":    users,
		"message": "you succefully got users",
	})
}

func (r *Repository) GetUserByUsername(c *gin.Context) {
	var user models.User
	userName := c.Param("user_name")

	if err := r.DB.First(&user, fmt.Sprintf("user_name = '%s'", userName)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting user by user_name",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    user,
		"message": "you succefully got user by user_name",
	})
}

func (r *Repository) UpdateUserByUsername(c *gin.Context) {
	var user models.User
	var foundUser models.User
	var count int64
	userNameKey := fmt.Sprint(c.Keys["user_name"])

	uid := fmt.Sprint(c.Keys["uid"])
	userName := c.Param("user_name")

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding json",
		})

		return
	}

	if err := r.DB.First(&foundUser, fmt.Sprintf("user_name = '%s'", userName)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding user",
		})

		return
	}

	if userNameKey != foundUser.UserName {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have permission to access this",
			"message": "error occured while checking user permission",
		})

		return
	}

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

	password := utils.HashPassword(user.Password)
	user.Password = password

	if foundUser.Email != user.Email {
		code := utils.Encode(randstr.String(20))
		user.EmailCode = &code
		user.EmailConfirmed = false

		emailData := utils.EmailData{
			URL:       r.Config.App_Origin + "/auth/confirm-email/" + *user.EmailCode,
			FirstName: user.FirstName,
			Subject:   "Your account verification code",
		}

		utils.SendEmail(&user, &emailData, r.Config)

		message := "you succefully update user. We sent an email with a verification code to " + user.Email

		if err := r.DB.Where("id = ?", uid).Updates(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"error":   err.Error(),
				"message": "error occured while update user",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"data":    user,
			"message": message,
		})

		return
	}

	if err := r.DB.Where("id = ?", uid).Updates(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while update user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    user,
		"message": "you succefully update user",
	})
}

func (r *Repository) DeleteUserByUsername(c *gin.Context) {
	var user models.User

	userNameKey := fmt.Sprint(c.Keys["user_name"])
	userName := c.Param("user_name")

	if userNameKey != userName {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have permission to access this",
			"message": "error occured while checking user permission",
		})

		return
	}

	if err := r.DB.First(&user, fmt.Sprintf("user_name = '%s'", userNameKey)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding user",
		})

		return
	}

	if err := r.DB.Delete(&user, fmt.Sprintf("user_name = %s", userNameKey)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while dekete user",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "you succefully delete user by username",
	})
}

func (r *Repository) AdminUpdateUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding json",
		})

		return
	}

	userNameKey := fmt.Sprint(c.Keys["user_name"])
	userRoleKey := fmt.Sprint(c.Keys["role"])
	userName := c.Param("user_name")

	if userNameKey == userName {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "you can't ban yourself",
			"message": "error occured while admin updating user",
		})

		return
	}

	if userRoleKey != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have permission to ban a user",
			"message": "error occured while admin updated user",
		})

		return
	}

	if err := r.DB.Where("user_name = ?", userName).Updates(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while updating user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    user,
		"message": "you succefully update post",
	})
}
