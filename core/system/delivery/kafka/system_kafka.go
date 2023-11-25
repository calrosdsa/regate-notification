package kafka

import (
	"context"
	"fmt"
	"log"
	r "notification/domain/repository"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type SystemHandler struct {
	systemU r.SystemUseCase
}

func NewKafkaHandler(systemU r.SystemUseCase) SystemHandler {
	return SystemHandler{
		systemU: systemU,
	}
}

func (k *SystemHandler) NotificationDiffusionConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{	
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "notification-diffusion",
		GroupID:   "consumer-group-diffusion",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		log.Println("RUNNN")
		fmt.Printf("message at offset %d: %s = %s\n %s", m.Offset, string(m.Key), string(m.Value), m.Time.Local().String())
	    k.systemU.SendNotificationDiffusion(context.TODO(), m.Value)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

