package usecase

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	r "notification/domain/repository"

	"github.com/goccy/go-json"

	// "strconv"

	// domain "notification/domain"
	"time"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
)

type grupoUcase struct {
	messageRepo r.GrupoRepository
	timeout     time.Duration
	firebase    *firebase.App
	utilU       r.UtilUseCase
}

func NewUseCase(messageRepo r.GrupoRepository, firebase *firebase.App, timeout time.Duration,
	utilU r.UtilUseCase) r.GrupoUseCase {
	return &grupoUcase{messageRepo: messageRepo, firebase: firebase, timeout: timeout,
	utilU: utilU,}
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
	// tokens := make([]string, len(fcm_tokens))
	for _, val := range fcm_tokens {
		// tokens = append(tokens, val.FcmToken)
		u.utilU.SendNotification(ctx, val.FcmToken, payload, r.NotificationSalaCreation,u.firebase)
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
		u.utilU.SendNotification(ctx, val.FcmToken, byteMessages, r.NotificationMessageGroup,u.firebase)
	}
	payloadData := struct {
		Ids  []int  `json:"ids"`
		Data []byte `json:"data"`
	}{
		Ids:  ids,
		Data: message,
	}
	posturl := fmt.Sprintf("%s/ws/publish/grupo/message/", viper.GetString("hosts.main"))
	// JSON body
	body, err := json.Marshal(payloadData)

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
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
