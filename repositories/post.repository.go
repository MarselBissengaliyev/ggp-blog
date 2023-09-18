package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func (r *Repository) GetPosts(c *gin.Context) {
	var posts []models.Post

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
	)).Find(&posts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting posts",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    posts,
		"message": "you succefully got posts",
	})
}

func (r *Repository) GetPostBySlug(c *gin.Context) {
	var post models.Post
	slug := c.Param("slug")

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	post.ViewsCount += 1
	if err := r.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while incrementing post views_count",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stauts":  "success",
		"data":    post,
		"message": "you succefully found post by slug",
	})
}

func (r *Repository) CreatePost(c *gin.Context) {
	var post models.Post
	var count int64

	uid, err := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while reading uid from session",
		})

		return
	}

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding post json",
		})

		return
	}

	post.UserId = uint(uid)
	post.ViewsCount = 0
	post.IsBanned = false

	slug := slug.Make(post.Title)
	post.Slug = slug

	r.DB.Model(&models.Post{}).Where("slug = ?", slug).Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "this slug already exists",
			"message": "error occured while creating post",
		})
		return
	}

	if err := r.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating post",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    post,
		"message": "you succefully created post",
	})
}

func (r *Repository) UpdatePostBySlug(c *gin.Context) {
	var post models.Post
	var foundPost models.Post

	slug := c.Param("slug")
	uid, err := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while reading uid from session",
		})

		return
	}

	if err := r.DB.First(&foundPost, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"stauts":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug and user_id",
		})

		return
	}

	if foundPost.UserId != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying user_id of post",
		})

		return
	}

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding post json",
		})

		return
	}

	post.IsBanned = foundPost.IsBanned
	post.Slug = foundPost.Slug

	if err := r.DB.Where("id = ?", foundPost.ID).Updates(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while updating post",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    post,
		"message": "you succefully update post",
	})
}

func (r *Repository) DeletePostBySlug(c *gin.Context) {
	var post models.Post
	var foundPost models.Post

	slug := c.Param("slug")
	uid, err := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while reading uid from session",
		})

		return
	}

	if err := r.DB.First(&foundPost, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"stauts":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug and user_id",
		})

		return
	}

	if foundPost.UserId != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying user_id of post",
		})

		return
	}

	if err := r.DB.Delete(&post, fmt.Sprintf("user_id = %d AND slug = '%s'", uid, slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete post",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "you succefully delete post by slug",
	})
}
