package usecase

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	r "notification/domain/repository"
	"time"

	firebase "firebase.google.com/go"
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
	var consumes []r.Consumo
	for _,val := range res{
		consume := r.Consumo{
			TypeEntity: r.ReservaSala,
			IdEnitity: message.EntityId,
			Message: "Reserva de un cupo en una sala",
			Amount: val.Amount,
			ProfileId: val.ProfileId,
		}
		consumes = append(consumes, consume)
	}
	u.billingRepo.AddConsume(ctx,consumes)
	if err != nil {
		return
	}
	data, err := json.Marshal(message)
	for _, val := range res {
		u.utilU.SendNotification(ctx, val.FcmToken, data,notification, u.firebase)
	}
	return
}
