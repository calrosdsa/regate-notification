package usecase

import (
	"context"
	"fmt"
	"log"
	r "notification/domain/repository"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)
type utilUseCase struct {
	timeout time.Duration
	utilRepo r.UtilRepository
}

func NewUseCase(timeout time.Duration,utilRepo r.UtilRepository) r.UtilUseCase{
	return &utilUseCase{
		timeout: timeout,
		utilRepo: utilRepo,
	}
}

func (u *utilUseCase)SendMessageToKafka(w *kafka.Writer,data interface{},key string){
	json, err := json.Marshal(data)
		if err != nil {
			log.Println("Fail to parse", err)
		}
		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(key),
				Value: json,
			},
		)
		if err != nil {
			u.LogError("SendMessageToKafka","util_usecase",err.Error())
		}
}

func (u *utilUseCase)GetProfileFcmToken(ctx context.Context,id int)(res string,err error){
	ctx,cancel := context.WithTimeout(ctx,u.timeout)
	defer cancel()
	res,err = u.utilRepo.GetProfileFcmToken(ctx,id)
	return
}

func (u *utilUseCase) SendNotification(ctx context.Context, tokens string, data []byte, notificationType r.NotificationType,firebase *firebase.App){
	client, err := firebase.Messaging(ctx)
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
			"payload":  string(data),
			"type":     strconv.Itoa(int(notificationType)),
			"priority": "high",
		},
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println("fail to send message",err)
	}
	log.Println("Successfully sent message:", response)
}


func (u *utilUseCase) SendNotificationToAdmin(ctx context.Context, tokens string,title string,payload string,firebase *firebase.App){
	client, err := firebase.Messaging(ctx)
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
			"title":    title,
			"payload":  payload,
		},
	}
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println("FAIL TO SEND",err)
	} else {
		log.Println("Successfully sent message:", response)
	}
}

func (u *utilUseCase) SendNotificationMessage(ctx context.Context, tokens string,data string, notificationType r.NotificationType,firebase *firebase.App){
	client, err := firebase.Messaging(ctx)
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
			"payload":  data,
			"type":     strconv.Itoa(int(notificationType)),
			"priority": "high",
		},
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println("FAIL TO SEND",err)
	} else {
		log.Println("Successfully sent message:", response)
	}
}




func (u *utilUseCase)LogError(method string,file string,err string){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		logrus.Error("error opening file: %v", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
    })
    ctx.Error(err)
}

func (u *utilUseCase)LogInfo(method string,file string,err string){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		logrus.Fatalf("error opening file: %v", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
    })
    ctx.Info(err)
}


func (u *utilUseCase)CustomLog(method string,file string,err string,payload map[string]interface{}){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		log.Println("error opening file", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
		"extra":payload,
    })
    ctx.Error(err)
// 	if u.logger != nil {
// 	ctx := u.logger.WithFields(logrus.Fields{
// 		"method": method,
// 		"file":file,
// 		"extra":payload,
//     })
//     ctx.Error(err)
// }
}

func (u *utilUseCase)LogFatal(method string,file string,err string,payload map[string]interface{}){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		log.Println("error opening file", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
    })
    ctx.Error(err)
}