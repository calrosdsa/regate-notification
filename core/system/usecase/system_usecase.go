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

type systemUCase struct {
	firebase *firebase.App
	timeout  time.Duration
	utilU    r.UtilUseCase
	systemRepo r.SystemRepository
}

func NewUseCase(firebase *firebase.App, timeout time.Duration, utilU r.UtilUseCase,systemRepo r.SystemRepository) r.SystemUseCase {
	return &systemUCase{
		timeout:  timeout,
		firebase: firebase,
		utilU:    utilU,
		systemRepo: systemRepo,
	}
}



func (u *systemUCase)SendNotificationDiffusion(ctx context.Context,d []byte) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	var data r.RequestNotificationDiffusion
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("DATA ---------",data)
	// fcm_token,err := u.utilU.GetProfileFcmToken(ctx,data.EntityId)
	// if err != nil{
	// 	log.Println(err)
	// }else{
	// 	u.utilU.SendNotification(ctx,fcm_token,d,r.NotificationBilling,u.firebase)
	// }


	
	// err = u.SendNotificationUsersSala(ctx,message,r.NotificationSalaHasBeenReserved)
	// if err != nil {
	// 	log.Println("ERROR",err)
	// }
}

