package core

import (
	"database/sql"
	"log"
	"os/signal"
	"syscall"
	"time"

	// "time"
	//Grupo
	_messageKafka "notification/core/grupo/delivery/kafka"
	_messageRepo "notification/core/grupo/repository"
	_messageUcase "notification/core/grupo/usecase"

	//Sala
	_salaKafka "notification/core/sala/delivery/kafka"
	_salaRepo "notification/core/sala/repository"
	_salaUcase "notification/core/sala/usecase"

	//Conversation
	_conversationKafka "notification/core/conversation/delivery/kafka"
	_conversationRepo "notification/core/conversation/repository"
	_conversationUcase "notification/core/conversation/usecase"

	//Billing
	_billingKafka "notification/core/billing/delivery/kafka"
	_billingRepo "notification/core/billing/repository"
	_billingUcase "notification/core/billing/usecase"

	//System
	_systemKafka "notification/core/system/delivery/kafka"
	_systemRepo "notification/core/system/repository"
	_systemUcase "notification/core/system/usecase"

	_utilRepo "notification/core/util/repository"
	_utilUcase "notification/core/util/usecase"

	"os"

	_firebase "firebase.google.com/go"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func Init(db *sql.DB, firebase *_firebase.App) {
	timeout := time.Duration(5) * time.Second

	wsAccountW := &kafka.Writer{
		Addr:     kafka.TCP(viper.GetString("kafka.host")),
		Topic:    "notify-ws",
		Balancer: &kafka.LeastBytes{},
	}

	utilR := _utilRepo.NewRepo(db)
	utilU := _utilUcase.NewUseCase(timeout, utilR)

	billingR := _billingRepo.NewRepository(db)
	billingU := _billingUcase.NewUseCase(firebase, timeout, utilU, billingR)
	billinKafka := _billingKafka.NewKafkaHandler(billingU)

	systemR := _systemRepo.NewRepository(db)
	systemU := _systemUcase.NewUseCase(firebase, timeout, utilU, systemR)
	systeKafka := _systemKafka.NewKafkaHandler(systemU)

	grupoRepo := _messageRepo.NewRepository(db)
	grupoUcase := _messageUcase.NewUseCase(grupoRepo, firebase, timeout, utilU,wsAccountW)

	grupoKafka := _messageKafka.NewKafkaHandler(grupoUcase)

	salaRepo := _salaRepo.NewRepository(db)
	salaUseCase := _salaUcase.NewUseCase(salaRepo, firebase, timeout, utilU, billingR,wsAccountW)
	salaKafka := _salaKafka.NewKafkaHandler(salaUseCase)

	conversation := _conversationRepo.NewRepository(db)
	conversationU := _conversationUcase.NewUseCase(conversation, firebase, timeout, utilU,wsAccountW)
	conversationKafka := _conversationKafka.NewKafkaHandler(conversationU)

	go salaKafka.SalaReservationConflictConsumer()
	go grupoKafka.MessageGroupConsumer()
	go grupoKafka.SalaCreationConsumer()
	go salaKafka.SalaConsumer()
	go billinKafka.BillingNotificationConsumer()
	go salaKafka.MessageSalaConsumer()
	go conversationKafka.MessageConversationConsumer()
	go systeKafka.NotificationDiffusionConsumer()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	log.Println("Adios!")

}

// func forever() {
//     for {
//         log.Printf("%v+\n", time.Now())
//         time.Sleep(time.Second)
//     }
// }
