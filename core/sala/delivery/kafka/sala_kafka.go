package kafka

import (
	"context"
	// "fmt"
	"log"
	r "notification/domain/repository"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type SalaKafkaHander struct {
	salaU r.SalaUseCase
}

func NewKafkaHandler(salaU r.SalaUseCase) SalaKafkaHander {
	return SalaKafkaHander{
		salaU: salaU,
	}
}

func (k *SalaKafkaHander) SalaReservationConflictConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "sala-reservation-conflict",
		GroupID:   "consumer-group-sala-reservation-conflict",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		// log.Println("RUNNN")
		// fmt.Printf("message at offset %d: %s = %s\n %s", m.Offset, string(m.Key), string(m.Value), m.Time.Local().String())
		err = k.salaU.SalaReservationConflict(context.TODO(), m.Value)
		log.Println(err)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func (k *SalaKafkaHander) SalaConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{viper.GetString("kafka.host")},
		Topic:     "sala",
		GroupID:   "consumer-group-sala",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		// log.Println("RUNNN")
		// fmt.Printf("message at offset %d: %s = %s\n %s", m.Offset, string(m.Key), string(m.Value), m.Time.Local().String())
		err = k.salaU.SalaSendNotification(context.TODO(), m.Value)
		log.Println(err)
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
