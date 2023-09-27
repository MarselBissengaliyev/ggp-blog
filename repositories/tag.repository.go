package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/gin-gonic/gin"
)

func (r *Repository) GetTags(c *gin.Context) {
	var tags []models.Tag
	var post models.Post
	slug := c.Param("slug")

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "success",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if err := r.DB.Find(
		&tags,
		fmt.Sprintf("post_id = %d", post.ID),
	).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "success",
			"error":   err.Error(),
			"message": "error occured while finding tags by post_id",
		})

		return
	}

	var result []gin.H

	for _, tag := range tags {
		result = append(result, gin.H{
			"id":        tag.ID,
			"name":      tag.Name,
			"post_slug": slug,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"message": "you succesully got tags by id",
	})
}

func (r *Repository) UpdateTags(c *gin.Context) {
	var post models.Post
	var tags []models.Tag
	slug := c.Param("slug")
	uid, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := c.BindJSON(&tags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding tag json",
		})

		return
	}

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if post.UserId != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying author of post",
		})

		return
	}

	if err := r.DB.Where("post_id = ?", post.ID).Updates(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while updating tags",
		})

		return
	}

	var result []gin.H

	for _, tag := range tags {
		result = append(result, gin.H{
			"id":        tag.ID,
			"name":      tag.Name,
			"post_slug": slug,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"message": "you succefully update tags",
	})
}

func (r *Repository) DeleteTag(c *gin.Context) {
	var post models.Post
	var tags []models.Tag
	slug := c.Param("slug")
	uid, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if post.UserId != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying author of post",
		})

		return
	}

	if err := r.DB.Where("post_id = ?", post.ID).Delete(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete tags",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "you succefully delete tags",
	})
}
