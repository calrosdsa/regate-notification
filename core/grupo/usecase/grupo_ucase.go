package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notification/domain/repository"
	"strconv"

	// domain "notification/domain"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/spf13/viper"
)

type grupoUcase struct {
	messageRepo repository.GrupoRepository
	timeout     time.Duration
	firebase    *firebase.App
}

func NewUseCase(messageRepo repository.GrupoRepository, firebase *firebase.App,timeout time.Duration) repository.GrupoUseCase {
	return &grupoUcase{messageRepo: messageRepo, firebase: firebase, timeout: timeout}
}

func(u *grupoUcase)SendNotificationSalaCreation(ctx context.Context,payload []byte)(err error){
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data repository.SalaPayload
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
		u.sendNotifications(ctx, val.FcmToken, payload,repository.NotificationSalaCreation)
	}
	return
}

func (u *grupoUcase) SendNotificationToUsersGroup(ctx context.Context, message []byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data repository.MessageGroupPayload
	err = json.Unmarshal(message, &data)
	if err != nil {
		return
	}
	messages, err := u.messageRepo.GetLastMessagesFromGroup(ctx, data.GrupoId)
	if err != nil {
		return
	}
	byteMessages, err := json.Marshal(messages)
	if err != nil {
		log.Println(byteMessages)
	}
	fcm_tokens, err := u.messageRepo.GetUsersFromGroup(ctx, data.GrupoId)
	if err != nil {
		return
	}
	log.Println(string(byteMessages))
	// tokens := make([]string, len(fcm_tokens))
	ids := make([]int,0)
	for _, val := range fcm_tokens {
		ids = append(ids, val.ProfileId)
		// tokens = append(tokens, val.FcmToken)
		u.sendNotifications(ctx, val.FcmToken, byteMessages,repository.NotificationMessageGroup)
	}
	payloadData := struct {
		Ids []int `json:"ids"`
		Data []byte `json:"data"`
	}{
		Ids: ids,
		Data: message,
	}
	posturl := fmt.Sprintf("%s/ws/publish/grupo/message/",viper.GetString("hosts.main"))
    // JSON body
	body,err := json.Marshal(payloadData)

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
func (u *grupoUcase) sendNotifications(ctx context.Context, tokens string, messages []byte,notificationType repository.NotificationType ) {
	client, err := u.firebase.Messaging(ctx)
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
			"payload":  string(messages),
			"type":     strconv.Itoa(int(notificationType)),
			"priority": "high",
		},
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully sent message:", response)
}
