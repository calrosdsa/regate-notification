package repository

import (
	"context"

	firebase "firebase.google.com/go"
)

type UtilUseCase interface {
	SendNotification(ctx context.Context, tokens string, data []byte, notificationType NotificationType,firebase *firebase.App)
	SendNotificationMessage(ctx context.Context, tokens string,data string, notificationType NotificationType,firebase *firebase.App)
	SendNotificationToAdmin(ctx context.Context,tokens string,title string,payload string,firebase *firebase.App)

	GetProfileFcmToken(ctx context.Context,id int)(string,error)
	LogError(method string, file string, err string)
	LogInfo(method string, file string, err string)
	CustomLog(method string, file string, err string, payload map[string]interface{})
	LogFatal(method string, file string, err string, payload map[string]interface{})
}

type UtilRepository interface {
	GetProfileFcmToken(ctx context.Context,id int)(string,error)
}