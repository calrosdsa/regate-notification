package kafka

import (
	"context"
	"fmt"
	"log"
	"notification/domain/repository"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type GrupoKafkaHandler struct {
	grupoUcase repository.GrupoUseCase
}

func NewKafkaHandler(grupoUcase repository.GrupoUseCase) GrupoKafkaHandler {
	return GrupoKafkaHandler{
		grupoUcase: grupoUcase,
	}
}


func (k *GrupoKafkaHandler) SalaCreationConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "sala-creation-group",
		GroupID:   "consumer-group-sala-creation",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	// r.SetOffset(2)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Erro read",err)
			break
		}
		log.Println("RUNNN")
		fmt.Printf("message at offset %d: %s = %s\n %s", m.Offset, string(m.Key), string(m.Value), m.Time.Local().String())
		err = k.grupoUcase.SendNotificationSalaCreation(context.Background(), m.Value)
		log.Println(err)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}


func (k *GrupoKafkaHandler) MessageGroupConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "notification-message-group",
		GroupID:   "consumer-group-messages",
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
		err = k.grupoUcase.SendNotificationToUsersGroup(context.Background(), m.Value)
		log.Println(err)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
