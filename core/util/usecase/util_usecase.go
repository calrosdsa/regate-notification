package usecase

import (
	"context"
	"log"
	r "notification/domain/repository"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)
type utilUseCase struct {
	timeout time.Duration
	utilRepo r.UtilRepository
}

func NewUseCase(timeout time.Duration,utilRepo r.UtilRepository) r.UtilUseCase{
	return &utilUseCase{
		timeout: timeout,
		utilRepo: utilRepo,
	}
}

func (u *utilUseCase)GetProfileFcmToken(ctx context.Context,id int)(res string,err error){
	ctx,cancel := context.WithTimeout(ctx,u.timeout)
	defer cancel()
	res,err = u.utilRepo.GetProfileFcmToken(ctx,id)
	return
}

func (u *utilUseCase) SendNotification(ctx context.Context, tokens string, data []byte, notificationType r.NotificationType,firebase *firebase.App){
	client, err := firebase.Messaging(ctx)
	if err != nil {
		log.Println(err)
	}
	message := &messaging.Message{
		//Notification: &messaging.Notification{
		//	Title: "Notification Test",
		//	Body:  "Hello React!!",
		//},
		Token: tokens,
		Data: map[string]string{
			"title":    "Nuevo Mensaje",
			"payload":  string(data),
			"type":     strconv.Itoa(int(notificationType)),
			"priority": "high",
		},
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println("fail to send message",err)
	}
	log.Println("Successfully sent message:", response)
}

func (u *utilUseCase) SendNotificationMessage(ctx context.Context, tokens string,data string, notificationType r.NotificationType,firebase *firebase.App){
	client, err := firebase.Messaging(ctx)
	if err != nil {
		log.Println(err)
	}
	message := &messaging.Message{
		//Notification: &messaging.Notification{
		//	Title: "Notification Test",
		//	Body:  "Hello React!!",
		//},
		Token: tokens,
		Data: map[string]string{
			"title":    "Nuevo Mensaje",
			"payload":  data,
			"type":     strconv.Itoa(int(notificationType)),
			"priority": "high",
		},
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println("FAIL TO SEND",err)
	} else {
		log.Println("Successfully sent message:", response)
	}
}