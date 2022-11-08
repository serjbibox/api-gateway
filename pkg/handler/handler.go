package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/serjbibox/api-gateway/pkg/storage"
)

// Обработчик HTTP запросов сервера GoNews
type Handler struct {
	storage *storage.Storage
}

//Конструктор объекта Handler
func New(storage *storage.Storage) (*Handler, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	return &Handler{storage: storage}, nil
}

//Инициализация маршрутизатора запросов.
//Регистрация обработчиков запросов
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	posts := router.Group("/news")
	{
		posts.POST("/comments", h.addComment)
		posts.GET("/latest", h.getNews)
		posts.GET("/:id", h.getDetailed)
	}
	return router
}
