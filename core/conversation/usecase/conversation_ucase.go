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

type conversationUseCase struct {
	conversationRepo r.ConversationRepository
	timeout     time.Duration
	firebase    *firebase.App
	utilU       r.UtilUseCase
}

func NewUseCase(conversationRepo r.ConversationRepository, firebase *firebase.App,
	timeout time.Duration, utilU r.UtilUseCase) r.ConversationUseCase {
	return &conversationUseCase{
		conversationRepo: conversationRepo,
		firebase:    firebase,
		timeout:     timeout,
		utilU:       utilU,
	}
}

func (u *conversationUseCase) SendNotificationMessageConversation(ctx context.Context,
	message []byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.Message
	err = json.Unmarshal(message, &data)
	if err != nil {
		return
	}

	// log.Println("GRUPOID", data.ParentId)
	// fcm_tokens, err := u.conversationRepo.(ctx, data.ParentId)
	// if err != nil {
	// 	return
	// }
	// tokens := make([]string, len(fcm_tokens))
	ids := make([]int, 0)
	if data.IsUser {
		fcm_tokens,err := u.conversationRepo.GetFcmTokenFromEstablecimientoAdmins(ctx,data.ParentId)
		if err != nil {
			log.Println("TOKEN ERROR", err)
		}
		for _,val := range fcm_tokens{
			log.Println("TOKEN",val.FcmToken)
			if val.FcmToken != nil {
				u.utilU.SendNotificationToAdmin(ctx,*val.FcmToken,"Tienes un nuevo mensaje",data.Content,u.firebase)
			}
			ids = append(ids, val.ProfileId)
		}
		
	} else {
		log.Println(string(message))
		messages, err1 := u.conversationRepo.GetLastMessagesConversation(ctx, data.ChatId)
		if err1 != nil {
			return
		}
		byteMessages, err := json.Marshal(messages)
		if err != nil {
			log.Println(byteMessages)
		}
		log.Println(data.ParentId)
		token, err := u.utilU.GetProfileFcmToken(ctx, data.ProfileId)
		if err != nil {
			log.Println("TOKEN ERROR", err)
		}
		u.utilU.SendNotification(ctx, token, byteMessages, r.NotificationMessageGroup, u.firebase)
		// u.utilU.SendNotification(ctx, "ci1jroyqZQbIL6ff4HP5nt:APA91bEildt9xtLetygsevwnf67SrultNKu-zhUXSd1LotU2VCLrDXlmHb_l_ndrAb4Mu554dX1EdF5D0o5dDwni_Mthf2Q3O8AHocmyilyar4enB7ATc9W2KhuhPtAAvj9BGgWES9qd",
			// byteMessages, r.NotificationMessageGroup, u.firebase)
	}
	// }
	if !data.IsUser {

		ids = append(ids, data.ProfileId)
		payloadData := r.MessageNotify{
		Ids:      ids,
		Data:     message,
		SenderId: 0,
		}
		posturl := fmt.Sprintf("%s/ws/publish/grupo/message/", viper.GetString("hosts.main"))
		// JSON body
		body, err := json.Marshal(payloadData)
		if err != nil {
			log.Println(err)
		}
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
	}else {

		payloadData := r.MessageNotify{
		Data:     message,
		Ids: ids,
		}
		posturl := fmt.Sprintf("%s/ws/publish/user/admin/", viper.GetString("hosts.main"))
		// JSON body
		body, err := json.Marshal(payloadData)
		if err != nil {
			log.Println(err)
		}
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
	}
	// log.Println("TOKENS", tokens)

	return
}