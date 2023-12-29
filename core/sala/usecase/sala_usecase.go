package usecase

import (
	"context"
	"encoding/json"

	// "fmt"
	"log"
	r "notification/domain/repository"
	"time"

	firebase "firebase.google.com/go"
	"github.com/segmentio/kafka-go"
)

type salaUseCase struct {
	salaRepo r.SalaRepository
	firebase *firebase.App
	timeout  time.Duration
	utilU    r.UtilUseCase
	billingRepo r.BillingRepository
	wsAccountW  *kafka.Writer

}

func NewUseCase(salaRepo r.SalaRepository, firebase *firebase.App, timeout time.Duration, utilU r.UtilUseCase,
	billingRepo r.BillingRepository,wsAccountW *kafka.Writer) r.SalaUseCase {
	return &salaUseCase{
		salaRepo: salaRepo,
		timeout:  timeout,
		firebase: firebase,
		utilU:    utilU,
		billingRepo: billingRepo,
		wsAccountW: wsAccountW,
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
		if val.FcmToken != nil {
			u.utilU.SendNotification(ctx, *val.FcmToken, byteMessages, r.NotificationMessageGroup,u.firebase)
		}
	}
	payloadData := r.MessageNotify{
		Ids:  ids,
		Data: message,
		SenderId: data.ProfileId,
	}
	u.utilU.SendMessageToKafka(u.wsAccountW,payloadData,string(r.KafkaPublishMessageEvent))
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
		if val.FcmToken != nil {
			u.utilU.SendNotification(ctx, *val.FcmToken, data,notification, u.firebase)
		}
	}
	return
}
func (u *salaUseCase) SendNotificationUsersSala(ctx context.Context,message r.MessageNotification,
	notification  r.NotificationType) (err error) {
	res, err := u.salaRepo.GetFcmTokensUserSalasSala(ctx, message.EntityId)
	data, err := json.Marshal(message)
	for _, val := range res {
		if val.FcmToken != nil{
			u.utilU.SendNotification(ctx,*val.FcmToken, data,notification, u.firebase)
		}
	}
	return
}
