package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/gin-gonic/gin"
)

func (r *Repository) GetComments(c *gin.Context) {
	var comments []models.Comment
	var post models.Post

	slug := c.Param("slug")

	limit := 10

	page, err := strconv.Atoi(c.Query("page"))
	offset := limit * page

	if err != nil || page == 1 {
		page = 1
		offset = 0
	}

	orderBy := "created_at"
	orderType := "desc"

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if err := r.DB.Limit(limit).Offset(offset).Order(fmt.Sprintf(
		"%s %s",
		orderBy,
		orderType,
	)).Preload("User").Find(&comments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting comments",
		})

		return
	}

	var result []gin.H

	// Loop through the fetched posts and get the count of PostReactions for each post
	for _, comment := range comments {
		result = append(result, gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"author":     comment.User.UserName,
			"post_slug":  slug,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"message": "you succefully got posts",
	})
}

func (r *Repository) CreateComment(c *gin.Context) {
	var comment models.Comment
	var post models.Post

	slug := c.Param("slug")
	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))
	userNameKey := fmt.Sprint(c.Keys["user_name"])

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding comment json",
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

	comment.PostId = post.ID
	comment.UserId = uint(userIdKey)

	if err := r.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating comment",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"author":     userNameKey,
			"post_slug":  slug,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
		},
		"message": "you succefully created a new comment",
	})
}

func (r *Repository) UpdateComment(c *gin.Context) {
	var comment models.Comment
	var foundComment models.Comment
	var post models.Post

	slug := c.Param("slug")
	id := c.Param("id")
	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))
	userNameKey := fmt.Sprint(c.Keys["user_name"])

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding comment json",
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

	if err := r.DB.First(&foundComment, fmt.Sprintf("id = '%s'", id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding comment by id",
		})

		return
	}

	if foundComment.UserId != uint(userIdKey) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this comment",
			"message": "error occured while verifying user_id of comment",
		})

		return
	}

	foundComment.Content = comment.Content

	if err := r.DB.Save(&foundComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while updating reaction",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stauts": "success",
		"data": gin.H{
			"id":         foundComment.ID,
			"content":    foundComment.Content,
			"author":     userNameKey,
			"post_slug":  slug,
			"created_at": foundComment.CreatedAt,
			"updated_at": foundComment.UpdatedAt,
		},
		"message": "you succefully update comment by id",
	})
}

func (r *Repository) DeleteComment(c *gin.Context) {
	var post models.Post
	var comment models.Comment

	slug := c.Param("slug")
	id := c.Param("id")
	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if err := r.DB.First(&comment, fmt.Sprintf("id = '%s'", id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding comment by id",
		})

		return
	}

	if comment.UserId != uint(userIdKey) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to delete this post-reaction",
			"message": "error occured while verifying author of post-reaction",
		})

		return
	}

	if err := r.DB.Delete(&comment, fmt.Sprintf("id = '%s'", id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete comment",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "you succefully delete comment by id",
	})
}
