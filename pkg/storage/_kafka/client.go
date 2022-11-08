package _kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Client — клиент Kafka.
type Client struct {
	// Reader осуществляет операции чтения из топика.
	Reader *kafka.Reader
	// Writer осуществляет операции записи в топик.
	Writer *kafka.Writer
}

// New создаёт и инициализирует клиента Kafka.
// Функция-конструктор.
func New(brokers []string, topic string, groupId int) (*Client, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	c := Client{}

	// Инициализация компонента получения сообщений.
	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: groupId,
		MinBytes:  10e1,
		MaxBytes:  10e6,
	})
	// Инициализация компонента отправки сообщений.
	c.Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &c, nil
}

// sendMessages отправляет сообщения в Kafka.
func (c *Client) SendMessages(messages []kafka.Message) error {
	err := c.Writer.WriteMessages(context.Background(), messages...)
	return err
}

// getMessage читает следующее сообщение из Kafka.
func (c *Client) GetMessage() (kafka.Message, error) {
	msg, err := c.Reader.ReadMessage(context.Background())
	return msg, err
}

// fetchProcessCommit сначала выбирает сообщение из очереди,
// потом обрабатывает, после чего подтверждает.
func (c *Client) FetchProcessCommit() (kafka.Message, error) {
	// Выборка очередного сообщения из Kafka.
	msg, err := c.Reader.FetchMessage(context.Background())
	if err != nil {
		return msg, err
	}

	// Обработка сообщения
	fmt.Println("обработка сообщения", string(msg.Key), string(msg.Value))

	// Подтверждение сообщения как обработанного.
	err = c.Reader.CommitMessages(context.Background(), msg)
	return msg, err
}
