package repositories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/gin-gonic/gin"
)

func (r *Repository) GetReactions(c *gin.Context) {
	var reactions []models.PostReaction
	var post models.Post

	postId := c.Param("post_id")

	limit := 20

	page, err := strconv.Atoi(c.Query("page"))
	offset := limit * page

	if err != nil || page == 1 {
		page = 1
		offset = 0
	}

	orderBy := "created_at"
	orderType := "desc"

	if err := r.DB.First(&post, fmt.Sprintf("id = '%s'", postId)).Error; err != nil {
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
	)).Find(&reactions).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while getting reactions",
		})

		return
	}

	var result []gin.H

	for _, reaction := range reactions {
		item := gin.H{
			"id":          reaction.ID,
			"is_liked":    reaction.IsLiked,
			"is_disliked": reaction.IsDisliked,
			"user_id":     reaction.UserId,
			"created_at":  reaction.CreatedAt,
			"updated_at":  reaction.UpdatedAt,
			"post_id":     reaction.PostId,
		}

		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"message": "you succefully got posts",
	})
}

func (r *Repository) CreateReaction(c *gin.Context) {
	var reaction models.PostReaction
	var post models.Post
	var count int64

	postId := c.Param("post_id")
	userId, err := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while converting uid from session",
		})

		return
	}

	if err := c.BindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding reaction json",
		})

		return
	}

	if err := r.DB.First(&post, fmt.Sprintf("id = '%s'", postId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	r.DB.Model(&models.PostReaction{}).Where("user_id = ?", userId).Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   "this reaction already exists",
			"message": "error occured while creating reaciton",
		})
		return
	}

	reaction.UserId = uint(userId)
	reaction.PostId = post.ID

	if err := r.DB.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while creating reaction",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"id":          reaction.ID,
			"is_liked":    reaction.IsLiked,
			"is_disliked": reaction.IsDisliked,
			"user_id":     reaction.UserId,
			"post_id":     reaction.PostId,
			"created_at":  reaction.CreatedAt,
			"updated_at":  reaction.UpdatedAt,
		},
		"message": "you succefully created a new reaction",
	})
}

func (r *Repository) UpdateReaction(c *gin.Context) {
	var reaction models.PostReaction
	var foundReaction models.PostReaction
	var post models.Post

	postId := c.Param("post_id")
	reactionId := c.Param("reaction_id")
	userId, err := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while converting uid from session",
		})

		return
	}

	if err := c.BindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while binding reaction json",
		})

		return
	}

	if err := r.DB.First(&post, fmt.Sprintf("id = '%s'", postId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if err := r.DB.First(&foundReaction, fmt.Sprintf("id = '%s'", reactionId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding reaction by id",
		})

		return
	}

	if foundReaction.UserId != uint(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to update this reaction",
			"message": "error occured while verifying user_id of reaction",
		})

		return
	}

	foundReaction.IsLiked = reaction.IsLiked
	foundReaction.IsDisliked = reaction.IsDisliked

	if err := r.DB.Save(&foundReaction).Error; err != nil {
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
			"id":          foundReaction.ID,
			"is_liked":    foundReaction.IsLiked,
			"is_disliked": foundReaction.IsDisliked,
			"user_id":     foundReaction.UserId,
			"post_id":     foundReaction.PostId,
			"created_at":  foundReaction.CreatedAt,
			"updated_at":  foundReaction.UpdatedAt,
		},
		"message": "you succefully update reaction by id",
	})
}

func (r *Repository) DeleteReaction(c *gin.Context) {
	var post models.Post
	var reaction models.PostReaction

	postId := c.Param("post_id")
	reactionId := c.Param("reaction_id")
	userId, err := strconv.Atoi(fmt.Sprint(c.Keys["uid"]))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while converting uid from session",
		})

		return
	}

	if err := r.DB.First(&post, fmt.Sprintf("id = '%s'", postId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post by slug",
		})

		return
	}

	if err := r.DB.First(&reaction, fmt.Sprintf("id = '%s'", reactionId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while finding post-reaction by id",
		})

		return
	}

	if reaction.UserId != uint(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"error":   "you don't have rights to delete this post-reaction",
			"message": "error occured while verifying author of post-reaction",
		})

		return
	}

	if err := r.DB.Delete(&reaction, fmt.Sprintf("id = '%s'", reactionId)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   err.Error(),
			"message": "error occured while delete post-reaction",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "you succefully delete post-reaction by id",
	})
}
