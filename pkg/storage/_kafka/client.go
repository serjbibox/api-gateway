package _kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaCache struct {
	Timestamp time.Time
	m         sync.Mutex
	Data      map[string]string
}

// Client — клиент Kafka.
type Client struct {
	// Reader осуществляет операции чтения из топика.
	Reader map[string]*kafka.Reader
	// Writer осуществляет операции записи в топик.
	Writer map[string]*kafka.Writer

	Buffer map[string]*KafkaCache
}

// New создаёт и инициализирует клиента Kafka.
// Функция-конструктор.
func New(ctx context.Context, brokers []string, topic map[string]string, groupId string) (*Client, error) {
	if len(brokers) == 0 || topic == nil {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}
	c := Client{
		Reader: make(map[string]*kafka.Reader),
		Writer: make(map[string]*kafka.Writer),
		Buffer: make(map[string]*KafkaCache),
	}
	for writeTopic, readTopic := range topic {
		c.Reader[readTopic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    readTopic,
			GroupID:  groupId,
			MinBytes: 10e1,
			MaxBytes: 10e6,
		})
		c.Writer[writeTopic] = &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    writeTopic,
			Balancer: &kafka.LeastBytes{},
		}
		c.Buffer[readTopic] = &KafkaCache{
			Data: make(map[string]string),
		}

	}
	go reader(ctx, &c)
	return &c, nil
}

func reader(ctx context.Context, c *Client) {
	for {
		m, err := c.GetMessage(ctx)
		if err != nil {
			log.Printf("failed reading from %s topic with error %w", c.Reader.Config().Topic, err)
		}
		log.Printf("Reader: %s:%s", m.Key, m.Value)
		c.Buffer.m.Lock()
		c.Buffer.Data[string(m.Key)] = string(m.Value)
		c.Buffer.m.Unlock()
	}
}

// sendMessages отправляет сообщения в Kafka.
func (c *Client) SendMessages(ctx context.Context, messages []kafka.Message) error {
	err := c.Writer.WriteMessages(ctx, messages...)
	return err
}

// getMessage читает следующее сообщение из Kafka.
func (c *Client) GetMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.Reader.ReadMessage(ctx)
	return msg, err
}

// fetchProcessCommit сначала выбирает сообщение из очереди,
// потом обрабатывает, после чего подтверждает.
func (c *Client) FetchProcessCommit(ctx context.Context) (kafka.Message, error) {
	// Выборка очередного сообщения из Kafka.
	msg, err := c.Reader.FetchMessage(ctx)
	if err != nil {
		return msg, err
	}

	// Обработка сообщения
	fmt.Println("обработка сообщения", string(msg.Key), string(msg.Value))

	// Подтверждение сообщения как обработанного.
	err = c.Reader.CommitMessages(context.Background(), msg)
	return msg, err
}
