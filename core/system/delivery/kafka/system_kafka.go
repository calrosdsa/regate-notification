package kafka

import (
	"context"
	"fmt"
	"log"
	r "notification/domain/repository"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type BillingHandler struct {
	billingU r.BillingUseCase
}

func NewKafkaHandler(billingU r.BillingUseCase) BillingHandler {
	return BillingHandler{
		billingU: billingU,
	}
}

func (k *BillingHandler) BillingNotificationConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{	
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "system",
		GroupID:   "consumer-group-system",
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
	    k.billingU.SendNotificationUserBilling(context.TODO(), m.Value)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

