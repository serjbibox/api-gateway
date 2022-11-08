package storage

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/serjbibox/api-gateway/pkg/models"
	"github.com/serjbibox/api-gateway/pkg/storage/_kafka"
)

var elog = log.New(os.Stderr, "Storage error\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "Storage info\t", log.Ldate|log.Ltime)

//задаёт контракт на работу с таблицей публикаций БД.
type NewsItem interface {
	GetList() ([]models.NewsItem, error)            // получение всех публикаций
	GetDetailed(id string) (models.NewsItem, error) // Получение детализированной публикации
}

//задаёт контракт на работу с таблицей публикаций БД.
type CommentsItem interface {
	Create()
}

// Хранилище данных.
type Storage struct {
	NewsItem
}

// Конструктор объекта хранилища для БД NewsgreSQL.
func NewStorageKafka(ctx context.Context, db *_kafka.Client) (*Storage, error) {
	if ctx == nil {
		elog.Println("context is nil")
		return nil, errors.New("context is nil")
	}
	if db == nil {
		elog.Println("db is nil")
		return nil, errors.New("db is nil")
	}
	return &Storage{
		NewsItem: newNewsKafka(ctx, db),
	}, nil
}
