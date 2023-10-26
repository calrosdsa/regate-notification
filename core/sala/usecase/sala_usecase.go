package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "fmt"
	"log"
	r "notification/domain/repository"
	"time"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
)

type salaUseCase struct {
	salaRepo r.SalaRepository
	firebase *firebase.App
	timeout  time.Duration
	utilU    r.UtilUseCase
	billingRepo r.BillingRepository
}

func NewUseCase(salaRepo r.SalaRepository, firebase *firebase.App, timeout time.Duration, utilU r.UtilUseCase,
	billingRepo r.BillingRepository) r.SalaUseCase {
	return &salaUseCase{
		salaRepo: salaRepo,
		timeout:  timeout,
		firebase: firebase,
		utilU:    utilU,
		billingRepo: billingRepo,
	}
}
func (u *salaUseCase)SendNotificationMessageSala(ctx context.Context,message []byte)(err error){
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.Message
	err = json.Unmarshal(message, &data)
	if err != nil {
		return
	}
	log.Println(string(message))
	messages, err := u.salaRepo.GetLastMessagesFromSala(ctx, data.ChatId)
	if err != nil {
		return
	}
	byteMessages, err := json.Marshal(messages)
	if err != nil {
		log.Println(byteMessages)
	}
	log.Println("GRUPOID",data.ParentId)
	fcm_tokens, err := u.salaRepo.GetFcmTokensUserSalasSala(ctx, data.ParentId)
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


func (u *salaUseCase) SalaSendNotification(ctx context.Context,d []byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.MessageNotification
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err)
		return
	}
	// log.Println(message)
	err = u.SendNotificationUsersSala(ctx,data,r.NotificationSalaHasBeenReserved)
	if err != nil {
		log.Println("ERROR",err)
	}
	return
}

func (u *salaUseCase) SalaReservationConflict(ctx context.Context,d []byte) (err error) {
	var data r.SalaConflictData
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("IDDDDDD",data.SalaIds)
	for _, val := range data.SalaIds {
		log.Println("IDDDDDD",val.Id)
		// message := r.MessageNotification{
		// 	Message:  "Lamentamos informarte que alguien más ha reservado la cancha que habías seleccionado para la sala.",
		// 	EntityId: val.Id,
		// }
		horario,err := u.salaRepo.GetSalaReservaHora(ctx,val.Id)
		horario.Id = val.Id
		horario.Message = "¡No te quedes sin jugar! Tenemos más canchas disponibles."
		if err == nil {
			err = u.SendNotificationUsersSala2(ctx,horario,r.NotificationSalaReservationConflict)
			if err != nil {
				log.Println("ERROR",err)
			}
		}
		// err = u.salaRepo.DeleteSala(ctx,val.Id)
		// if err != nil{
		// 	log.Println()
		// 	return 
		// }
	}
	return
}
func (u *salaUseCase) SendNotificationUsersSala2(ctx context.Context,message r.SalaHora,
	notification  r.NotificationType) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res, err := u.salaRepo.GetFcmTokensUserSalasSala(ctx, message.Id)
	log.Println("SALA HORA-------",message)
	if err != nil {
		log.Println("FAILT TO FETCG TOKENS",err)
		return
	}
	data, err := json.Marshal(message)
	for _, val := range res {
		log.Println("FCM_TOKENS", val.FcmToken)
		u.utilU.SendNotification(ctx, val.FcmToken, data,notification, u.firebase)
	}
	return
}
func (u *salaUseCase) SendNotificationUsersSala(ctx context.Context,message r.MessageNotification,
	notification  r.NotificationType) (err error) {
	res, err := u.salaRepo.GetFcmTokensUserSalasSala(ctx, message.EntityId)
	data, err := json.Marshal(message)
	for _, val := range res {
		u.utilU.SendNotification(ctx, val.FcmToken, data,notification, u.firebase)
	}
	return
}
