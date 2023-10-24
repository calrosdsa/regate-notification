package repository

import (
	"context"

	firebase "firebase.google.com/go"
)

type UtilUseCase interface {
	SendNotification(ctx context.Context, tokens string, data []byte, notificationType NotificationType,firebase *firebase.App)
	SendNotificationMessage(ctx context.Context, tokens string,data string, notificationType NotificationType,firebase *firebase.App)
	GetProfileFcmToken(ctx context.Context,id int)(string,error)
}

type UtilRepository interface {
	GetProfileFcmToken(ctx context.Context,id int)(string,error)
}