package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/gin-gonic/gin"
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

	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "user_name param cannot be empty string",
			"message": "error occured while reading user_name param",
		})

		return
	}

	if err := r.DB.First(&user, fmt.Sprintf("user_name = %s", userName)).Error; err != nil {
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
