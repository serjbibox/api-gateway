package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serjbibox/api-gateway/pkg/models"
)

// Получение всех публикаций.
func (h *Handler) getNews(c *gin.Context) {
	posts, err := h.storage.NewsItem.GetList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}

// Обновление публикации.
func (h *Handler) getDetailed(c *gin.Context) {
	var post models.NewsItem
	if err := c.BindJSON(&post); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	news, err := h.storage.NewsItem.GetDetailed("1")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": news,
	})
}

// Получение всех публикаций.
func (h *Handler) addComment(c *gin.Context) {
	posts, err := h.storage.NewsItem.GetList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}
