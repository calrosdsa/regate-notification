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

type billingUCase struct {
	firebase *firebase.App
	timeout  time.Duration
	utilU    r.UtilUseCase
	billingRepo r.BillingRepository
}

func NewUseCase(firebase *firebase.App, timeout time.Duration, utilU r.UtilUseCase,billingRepo r.BillingRepository) r.BillingUseCase {
	return &billingUCase{
		timeout:  timeout,
		firebase: firebase,
		utilU:    utilU,
		billingRepo: billingRepo,
	}
}

// func (u *billingUCase) AddConsume(ctx context.Context,d r.Consumo) {
// 	ctx, cancel := context.WithTimeout(ctx, u.timeout)
// 	defer cancel()
// 	u.billingRepo.AddConsume(ctx,d)
// }

func (u *billingUCase) SendNotificationUserBilling(ctx context.Context,d []byte) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.MessageNotification
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("DATA ---------",data)
	fcm_token,err := u.utilU.GetProfileFcmToken(ctx,data.EntityId)
	if err != nil{
		log.Println("FAIL TO SEND MESSAGE",err)
	}else{
		u.utilU.SendNotification(ctx,fcm_token,d,r.NotificationBilling,u.firebase)
	}
	// err = u.SendNotificationUsersSala(ctx,message,r.NotificationSalaHasBeenReserved)
	// if err != nil {
	// 	log.Println("ERROR",err)
	// }
}

// func (u *billingUCase) SalaReservationConflict(ctx context.Context,d []byte) (err error) {
// 	ctx, cancel := context.WithTimeout(ctx, u.timeout)
// 	defer cancel()
// 	var data r.SalaConflictData
// 	err = json.Unmarshal(d, &data)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	log.Println("IDDDDDD",data.SalaIds)
// 	for _, val := range data.SalaIds {
// 		log.Println("IDDDDDD",val.Id)
// 		message := r.MessageNotification{
// 			Message:  "Lamentamos informarte que alguien más ha reservado la cancha que habías seleccionado para la sala.",
// 			EntityId: val.Id,
// 		}
// 		err = u.SendNotificationUsersSala(ctx,message,r.NotificationSalaReservationConflict)
// 		if err != nil {
// 			log.Println("ERROR",err)
// 		}
// 		err = u.salaRepo.DeleteSala(ctx,val.Id)
// 		if err != nil{
// 			log.Println()
// 			return 
// 		}
// 	}
// 	return
// }
// func (u *billingUCase) SendNotificationUsersSala(ctx context.Context,message r.MessageNotification,
// 	notification  r.NotificationType) (err error) {
// 	ctx, cancel := context.WithTimeout(ctx, u.timeout)
// 	defer cancel()
// 	res, err := u.salaRepo.GetFcmTokensUserSalasSala(ctx, message.EntityId)
// 	if err != nil {
// 		log.Println("FAILT TO FETCG TOKENS",err)
// 		return
// 	}
// 	data, err := json.Marshal(message)
// 	for _, val := range res {
// 		log.Println("FCM_TOKENS", val.FcmToken)
// 		u.utilU.SendNotification(ctx, val.FcmToken, data,notification, u.firebase)
// 	}
// 	return
// }
