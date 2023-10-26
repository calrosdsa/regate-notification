package kafka

import (
	"context"
	"fmt"
	"log"
	"notification/domain/repository"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type ConversationKafkaHandler struct {
	conversationU repository.ConversationUseCase
}

func NewKafkaHandler(conversationU repository.ConversationUseCase) ConversationKafkaHandler {
	return ConversationKafkaHandler{
		conversationU: conversationU,
	}
}

func (k *ConversationKafkaHandler) MessageConversationConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "notification-message-conversation",
		GroupID:   "consumer-conversation-messages",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	// r.SetOffset(2)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		log.Println("RUNNN")
		fmt.Printf("message at offset %d: %s = %s\n %s", m.Offset, string(m.Key), string(m.Value), m.Time.Local().String())
		err = k.conversationU.SendNotificationMessageConversation(context.Background(), m.Value)
		log.Println(err)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
