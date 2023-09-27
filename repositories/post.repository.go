package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/MarselBissengaliyev/ggp-blog/utils"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (r *Repository) GetPosts(c *gin.Context) {
	var posts []models.Post

	validOrderByValues := []string{"views_count", "created_at", "updated_at"}
	validOrderTypeValues := []string{"desc", "asc"}
	arrayUtil := new(utils.ArrayUtil)

	limit := 5

	page, err := strconv.Atoi(c.Query("page"))
	offset := limit * page

	if err != nil || page == 1 {
		page = 1
		offset = 0
	}

	orderBy := c.Query("order_by")

	if !arrayUtil.IsValid(orderBy, validOrderByValues) {
		orderBy = "views_count"
	}

	orderType := c.Query("order_type")

	if !arrayUtil.IsValid(orderType, validOrderTypeValues) {
		orderType = "desc"
	}

	if err := r.DB.Limit(limit).Offset(offset).Order(fmt.Sprintf(
		"%s %s",
		orderBy,
		orderType,
	)).Preload("User").Preload("PostReactions").Preload("Comments").Find(
		&posts,
		"is_banned = false",
	).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting posts",
		})

		return
	}

	var result []gin.H

	// Loop through the fetched posts and get the count of PostReactions for each post
	for i, post := range posts {
		var reactionCount int64 = r.DB.Model(&posts[i]).Association("PostReactions").Count()
		var commentCount int64 = r.DB.Model(&posts[i]).Association("Comments").Count()

		result = append(result, gin.H{
			"title":           post.Title,
			"description":     post.Description,
			"slug":            post.Slug,
			"views_count":     post.ViewsCount,
			"reactions_count": reactionCount,
			"comments_count":  commentCount,
			"content":         post.PreviewUrl,
			"author":          post.User.UserName,
			"tags":            post.Tags,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"message": "you succefully got posts",
	})
}

func (r *Repository) GetPostBySlug(c *gin.Context) {
	var post models.Post
	var reactionsCount int64
	var commentCount int64

	slug := c.Param("slug")

	if err := r.DB.Preload("User").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Model(&models.Comment{}).Count(&commentCount)
	}).Preload("PostReactions", func(db *gorm.DB) *gorm.DB {
		return db.Model(&models.PostReaction{}).Count(&reactionsCount)
	}).First(
		&post,
		fmt.Sprintf("slug = '%s' AND is_banned= false", slug),
	).Error; err != nil {
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
		"stauts": "success",
		"data": gin.H{
			"title":           post.Title,
			"description":     post.Description,
			"slug":            post.Slug,
			"views_count":     post.ViewsCount,
			"reactions_count": reactionsCount,
			"comments_count":  commentCount,
			"content":         post.PreviewUrl,
			"author":          post.User.UserName,
		},
		"message": "you succefully found post by slug",
	})
}

func (r *Repository) CreatePost(c *gin.Context) {
	var post models.Post
	var count int64

	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))
	userNameKey := fmt.Sprint(c.Keys["user_name"])

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding post json",
		})

		return
	}

	post.ViewsCount = 0
	post.UserId = uint(userIdKey)

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

	var tags []gin.H
	for _, tag := range post.Tags {
		tags = append(tags, gin.H{
			"name":      tag.Name,
			"post_slug": slug,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"title":       post.Title,
			"slug":        post.Slug,
			"description": post.Description,
			"content":     post.Content,
			"preview_url": post.PreviewUrl,
			"author":      userNameKey,
			"views_count": post.ViewsCount,
			"tags":        tags,
		},
		"message": "you succefully created post",
	})
}

func (r *Repository) UpdatePostBySlug(c *gin.Context) {
	var post models.Post
	var foundPost models.Post
	var reactionsCount int64
	var commentCount int64

	slug := c.Param("slug")
	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := r.DB.Preload("User").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Model(&models.Comment{}).Count(&commentCount)
	}).Preload("PostReactions", func(db *gorm.DB) *gorm.DB {
		return db.Model(&models.PostReaction{}).Count(&reactionsCount)
	}).First(&foundPost, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"stauts":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug and user_id",
		})

		return
	}

	if foundPost.UserId != uint(userIdKey) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying author of post",
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
		"status": "success",
		"data": gin.H{
			"title":           post.Title,
			"description":     post.Description,
			"slug":            post.Slug,
			"views_count":     foundPost.ViewsCount,
			"reactions_count": reactionsCount,
			"comments_count":  commentCount,
			"content":         post.PreviewUrl,
			"author":          foundPost.User.UserName,
		},
		"message": "you succefully update post",
	})
}

func (r *Repository) DeletePostBySlug(c *gin.Context) {
	var post models.Post
	var foundPost models.Post

	slug := c.Param("slug")
	userIdKey, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := r.DB.First(&foundPost, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"stauts":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug and user_id",
		})

		return
	}

	if foundPost.UserId != uint(userIdKey) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying user_id of post",
		})

		return
	}

	if err := r.DB.Delete(&post, fmt.Sprintf("user_id = %d AND slug = '%s'", userIdKey, slug)).Error; err != nil {
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
