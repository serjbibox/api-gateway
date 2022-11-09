package storage

import (
	"context"
	"log"

	"github.com/serjbibox/api-gateway/pkg/models"
	"github.com/serjbibox/api-gateway/pkg/storage/_kafka"
)

//Объект, реализующий интерфейс работы с таблицей posts NewsgreSQL.
type NewsKafka struct {
	db  *_kafka.Client
	ctx context.Context
}

//Конструктор NewsNewsgres
func newNewsKafka(ctx context.Context, db *_kafka.Client) NewsItem {
	return &NewsKafka{
		db:  db,
		ctx: ctx,
	}
}

func (p *NewsKafka) GetList() ([]models.NewsItem, error) {
	log.Println("запрос на чтение в кафку")
	m, err := p.db.GetMessage(p.ctx)
	log.Println("дождались")
	if err != nil {
		return nil, err
	}

	return []models.NewsItem{
		{
			Description: string(m.Value),
		},
	}, nil
}

func (p *NewsKafka) GetDetailed(id string) (models.NewsItem, error) {

	return models.NewsItem{}, nil
}
