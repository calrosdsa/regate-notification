package usecase

import (
	"context"
	"github.com/goccy/go-json"

	// "encoding/json"
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
	payload,err := json.Marshal(data.Notification)
	if err != nil{
		u.utilU.LogError("SendNotificationDiffusion","system_usecase",err.Error())
	}
	fcm_tokens,err := u.systemRepo.GetUserFcmTokens(ctx,data.Categories)
	if err != nil{
		u.utilU.LogError("SendNotificationDiffusion","system_usecase",err.Error())
	}
	for _,fcm :=range fcm_tokens {
		if fcm.FcmToken != nil {	
			u.utilU.SendNotification(ctx,*fcm.FcmToken,payload,r.NotificationEvent,u.firebase)
		}
	}
	

	
	// err = u.SendNotificationUsersSala(ctx,message,r.NotificationSalaHasBeenReserved)
	// if err != nil {
	// 	log.Println("ERROR",err)
	// }
}

