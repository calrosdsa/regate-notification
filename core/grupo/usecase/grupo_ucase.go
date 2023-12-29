package usecase

import (
	"context"
	"log"
	r "notification/domain/repository"

	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"

	// "strconv"

	// domain "notification/domain"
	"time"

	firebase "firebase.google.com/go"
)

type grupoUcase struct {
	messageRepo r.GrupoRepository
	timeout     time.Duration
	firebase    *firebase.App
	utilU       r.UtilUseCase
	wsAccountW  *kafka.Writer

}

func NewUseCase(messageRepo r.GrupoRepository, firebase *firebase.App, timeout time.Duration,
	utilU r.UtilUseCase,wsAccountW *kafka.Writer) r.GrupoUseCase {
	return &grupoUcase{
		messageRepo: messageRepo,
		firebase: firebase,
		timeout: timeout,
	    utilU: utilU,
		wsAccountW: wsAccountW,
	}
}

func (u *grupoUcase) SendNotificationSalaCreation(ctx context.Context, payload []byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.SalaPayload
	err = json.Unmarshal(payload, &data)
	if err != nil {
		return
	}
	fcm_tokens, err := u.messageRepo.GetUsersFromGroup(ctx, data.GrupoId)
	if err != nil {
		return
	}
	log.Println("SENDIUNG TO USERS SALA")
	// tokens := make([]string, len(fcm_tokens))
	for _, val := range fcm_tokens {
		// tokens = append(tokens, val.FcmToken)\
		if val.FcmToken != nil {
			if data.SenderId != val.ProfileId {
				u.utilU.SendNotification(ctx, *val.FcmToken, payload, r.NotificationSalaCreation,u.firebase)
			}
		}
	}
	return
}

func (u *grupoUcase) SendNotificationToUsersGroup(ctx context.Context, message []byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.Message
	err = json.Unmarshal(message, &data)
	if err != nil {
		return
	}
	log.Println(string(message))
	messages, err := u.messageRepo.GetLastMessagesFromGroup(ctx, data.ChatId)
	if err != nil {
		return
	}
	byteMessages, err := json.Marshal(messages)
	if err != nil {
		log.Println(byteMessages)
	}
	log.Println("GRUPOID", data.ParentId)
	fcm_tokens, err := u.messageRepo.GetUsersFromGroup(ctx, data.ParentId)
	if err != nil {
		return
	}
	log.Println(string(byteMessages))
	// tokens := make([]string, len(fcm_tokens))
	ids := make([]int, 0)
	for _, val := range fcm_tokens {
		ids = append(ids, val.ProfileId)
		// tokens = append(tokens, val.FcmToken)
		if val.FcmToken != nil {
			u.utilU.SendNotification(ctx, *val.FcmToken, byteMessages, r.NotificationMessageGroup,u.firebase)
		}
	}
	log.Println("NOTIFICATIONS SENDED")
	payloadData := r.MessageNotify{
		Ids:  ids,
		Data: message,
		SenderId: data.ProfileId,
	}
	u.utilU.SendMessageToKafka(u.wsAccountW,payloadData,string(r.KafkaPublishMessageEvent))
	// log.Println("TOKENS", tokens)
	return
}

// func (u *grupoUcase) sendNotifications(ctx context.Context, tokens string, messages []byte, notificationType repository.NotificationType) {
// 	client, err := u.firebase.Messaging(ctx)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	message := &messaging.Message{
// 		//Notification: &messaging.Notification{
// 		//	Title: "Notification Test",
// 		//	Body:  "Hello React!!",
// 		//},
// 		Token: tokens,
// 		Data: map[string]string{
// 			"title":    "Nuevo Mensaje",
// 			"payload":  string(messages),
// 			"type":     strconv.Itoa(int(notificationType)),
// 			"priority": "high",
// 		},
// 	}

// 	response, err := client.Send(ctx, message)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println("Successfully sent message:", response)
// }
