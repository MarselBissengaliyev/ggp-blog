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
	uid, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

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
		fmt.Sprintf("post_id = %d AND user_id = %d", post.ID, uint(uid)),
	).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "success",
			"error":   err.Error(),
			"message": "error occured while finding tags by post_id",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    tags,
		"message": "you succesully got tags by id",
	})
}

func (r *Repository) CreateTag(c *gin.Context) {
	var post models.Post
	var tag models.Tag
	var tagsCount int64
	slug := c.Param("slug")
	uid, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := c.BindJSON(&tag); err != nil {
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

	if post.ID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying author of post",
		})

		return
	}

	r.DB.Model(models.Tag{}).Where(
		"user_id = ? AND post_id = ? AND name = ?",
		uid, post.ID, tag.Name,
	).Count(&tagsCount)

	if tagsCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "this tag already exists",
			"message": "error occured while tag",
		})
		return
	}

	tag.PostId = post.ID

	if err := r.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating tag",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    tag,
		"message": "you succefully created tag",
	})
}

func (r *Repository) UpdateTag(c *gin.Context) {
	var post models.Post
	var tag models.Tag
	var foundTag models.Tag
	slug := c.Param("slug")
	tagId := c.Param("tag_id")
	uid, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := c.BindJSON(&tag); err != nil {
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

	if post.ID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying author of post",
		})

		return
	}

	if err := r.DB.First(&foundTag, fmt.Sprintf("id = '%s'", tagId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"stauts":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding tag by id",
		})

		return
	}

	foundTag.Name = tag.Name
	if err := r.DB.Save(&foundTag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while updating tag",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    foundTag,
		"message": "you succefully update tag",
	})
}

func (r *Repository) DeleteTag(c *gin.Context) {
	var post models.Post
	var tag models.Tag
	slug := c.Param("slug")
	tagId := c.Param("tag_id")
	uid, _ := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err := r.DB.First(&post, fmt.Sprintf("slug = '%s'", slug)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if post.ID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this post",
			"message": "error occured while verifying author of post",
		})

		return
	}

	if err := r.DB.First(&tag, fmt.Sprintf("id = '%s'", tagId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"stauts":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding tag by id",
		})

		return
	}

	if err := r.DB.Delete(&tag, fmt.Sprintf("id = '%s'", tagId)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete tag",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "you succefully delete tag by id",
	})
}
