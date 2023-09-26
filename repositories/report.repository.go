package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/gin-gonic/gin"
)

func (r *Repository) GetReports(c *gin.Context) {
	var reports []models.Report

	limit := 5

	page, err := strconv.Atoi(c.Query("page"))
	offset := limit * page

	if err != nil || page == 1 {
		page = 1
		offset = 0
	}

	orderBy := c.Query("order_by")

	if orderBy == "" {
		orderBy = "views_count"
	}

	orderType := c.Query("order_type")

	if orderType == "" {
		orderType = "desc"
	}

	if err := r.DB.Limit(limit).Offset(offset).Order(fmt.Sprintf(
		"%s %s",
		orderBy,
		orderType,
	)).Find(&reports).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting reports",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    reports,
		"message": "you succefully got reports",
	})
}

func (r *Repository) GetReportById(c *gin.Context) {
	var report models.Report
	reportId := c.Param("report_id")

	if err := r.DB.First(&report, fmt.Sprintf("id = '%s'", reportId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding report by id",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stauts":  "success",
		"data":    report,
		"message": "you succefully found post by slug",
	})
}

func (r *Repository) CreateReport(c *gin.Context) {
	var post models.Post
	var report models.Report
	var count int64

	slug := c.Param("slug")
	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if post.UserId == uint(userIdKey) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   "you can't report your own post",
			"message": "error occured checking post author",
		})

		return
	}

	report.PostId = post.ID
	report.UserId = uint(userIdKey)

	r.DB.Model(&models.Report{}).Where("user_id = ? AND post_id = ?", userIdKey, post.ID).Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "you have already sent the report",
			"message": "error occured while creating report",
		})
		return
	}

	if err := r.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating report",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    report,
		"message": "you succefully created report",
	})
}
